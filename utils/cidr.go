// Code taken and customized from https://raw.githubusercontent.com/projectdiscovery/mapcidr/24b047bfd50376c9e7f17123006d56cd6481d1e9/ip.go
package utils

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"net"
	"strings"
)

// Const for the CIDR calculation
const (
	ipv4BitLen = 8 * net.IPv4len
	ipv6BitLen = 8 * net.IPv6len
)

var (
	// v4Mappedv6Prefix is the RFC2765 IPv4-mapped address prefix.
	v4Mappedv6Prefix  = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff}
	ipv4LeadingZeroes = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	defaultIPv4       = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0}
	defaultIPv6       = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	upperIPv4         = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 255, 255, 255, 255}
	upperIPv6         = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

func ipNetToRange(ipNet net.IPNet) netWithRange {
	firstIP := make(net.IP, len(ipNet.IP))
	lastIP := make(net.IP, len(ipNet.IP))

	copy(firstIP, ipNet.IP)
	copy(lastIP, ipNet.IP)

	firstIP = firstIP.Mask(ipNet.Mask)
	lastIP = lastIP.Mask(ipNet.Mask)

	if firstIP.To4() != nil {
		firstIP = append(v4Mappedv6Prefix, firstIP...)
		lastIP = append(v4Mappedv6Prefix, lastIP...)
	}

	lastIPMask := make(net.IPMask, len(ipNet.Mask))
	copy(lastIPMask, ipNet.Mask)
	for i := range lastIPMask {
		lastIPMask[len(lastIPMask)-i-1] = ^lastIPMask[len(lastIPMask)-i-1]
		lastIP[net.IPv6len-i-1] |= lastIPMask[len(lastIPMask)-i-1]
	}

	return netWithRange{First: &firstIP, Last: &lastIP, Network: &ipNet}
}

func getPreviousIP(ip net.IP) net.IP {
	// Cannot go lower than zero!
	if ip.Equal(net.IP(defaultIPv4)) || ip.Equal(net.IP(defaultIPv6)) {
		return ip
	}

	previousIP := make(net.IP, len(ip))
	copy(previousIP, ip)

	var overflow bool
	var lowerByteBound int
	if ip.To4() != nil {
		lowerByteBound = net.IPv6len - net.IPv4len
	} else {
		lowerByteBound = 0
	}
	for i := len(ip) - 1; i >= lowerByteBound; i-- {
		if overflow || i == len(ip)-1 {
			previousIP[i]--
		}
		// Track if we have overflowed and thus need to continue subtracting.
		if ip[i] == 0 && previousIP[i] == 255 {
			overflow = true
		} else {
			overflow = false
		}
	}
	return previousIP
}

// GetNextIP returns the next IP from the given IP address. If the given IP is
// the last IP of a v4 or v6 range, the same IP is returned.
func GetNextIP(ip net.IP) net.IP {
	if ip.Equal(upperIPv4) || ip.Equal(upperIPv6) {
		return ip
	}

	nextIP := make(net.IP, len(ip))
	switch len(ip) {
	case net.IPv4len:
		ipU32 := binary.BigEndian.Uint32(ip)
		ipU32++
		binary.BigEndian.PutUint32(nextIP, ipU32)
		return nextIP
	case net.IPv6len:
		ipU64 := binary.BigEndian.Uint64(ip[net.IPv6len/2:])
		ipU64++
		binary.BigEndian.PutUint64(nextIP[net.IPv6len/2:], ipU64)
		if ipU64 == 0 {
			ipU64 = binary.BigEndian.Uint64(ip[:net.IPv6len/2])
			ipU64++
			binary.BigEndian.PutUint64(nextIP[:net.IPv6len/2], ipU64)
		} else {
			copy(nextIP[:net.IPv6len/2], ip[:net.IPv6len/2])
		}
		return nextIP
	default:
		return ip
	}
}

func createSpanningCIDR(r netWithRange) net.IPNet {
	// Don't want to modify the values of the provided range, so make copies.
	lowest := *r.First
	highest := *r.Last

	var isIPv4 bool
	var spanningMaskSize, bitLen, byteLen int
	if lowest.To4() != nil {
		isIPv4 = true
		bitLen = ipv4BitLen
		byteLen = net.IPv4len
	} else {
		bitLen = ipv6BitLen
		byteLen = net.IPv6len
	}

	if isIPv4 {
		spanningMaskSize = ipv4BitLen
	} else {
		spanningMaskSize = ipv6BitLen
	}

	// Convert to big Int so we can easily do bitshifting on the IP addresses,
	// since golang only provides up to 64-bit unsigned integers.
	lowestBig := big.NewInt(0).SetBytes(lowest)
	highestBig := big.NewInt(0).SetBytes(highest)

	// Starting from largest mask / smallest range possible, apply a mask one bit
	// larger in each iteration to the upper bound in the range  until we have
	// masked enough to pass the lower bound in the range. This
	// gives us the size of the prefix for the spanning CIDR to return as
	// well as the IP for the CIDR prefix of the spanning CIDR.
	for spanningMaskSize > 0 && lowestBig.Cmp(highestBig) < 0 {
		spanningMaskSize--
		mask := big.NewInt(1)
		mask = mask.Lsh(mask, uint(bitLen-spanningMaskSize))
		mask = mask.Mul(mask, big.NewInt(-1))
		highestBig = highestBig.And(highestBig, mask)
	}

	// If ipv4, need to append 0s because math.Big gets rid of preceding zeroes.
	if isIPv4 {
		highest = append(ipv4LeadingZeroes, highestBig.Bytes()...) //nolint
	} else {
		highest = highestBig.Bytes()
	}

	// Int does not store leading zeroes.
	if len(highest) == 0 {
		highest = make([]byte, byteLen)
	}

	newNet := net.IPNet{IP: highest, Mask: net.CIDRMask(spanningMaskSize, bitLen)}
	return newNet
}

type netWithRange struct {
	First   *net.IP
	Last    *net.IP
	Network *net.IPNet
}

// rangeToCIDRs converts the range of IPs covered by firstIP and lastIP to
// a list of CIDRs that contains all of the IPs covered by the range.
func rangeToCIDRs(firstIP, lastIP net.IP) []*net.IPNet {
	// First, create a CIDR that spans both IPs.
	spanningCIDR := createSpanningCIDR(netWithRange{&firstIP, &lastIP, nil})
	spanningRange := ipNetToRange(spanningCIDR)
	firstIPSpanning := spanningRange.First
	lastIPSpanning := spanningRange.Last

	cidrList := []*net.IPNet{}

	// If the first IP of the spanning CIDR passes the lower bound (firstIP),
	// we need to split the spanning CIDR and only take the IPs that are
	// greater than the value which we split on, as we do not want the lesser
	// values since they are less than the lower-bound (firstIP).
	if bytes.Compare(*firstIPSpanning, firstIP) < 0 {
		// Split on the previous IP of the first IP so that the right list of IPs
		// of the partition includes the firstIP.
		prevFirstRangeIP := getPreviousIP(firstIP)
		var bitLen int
		if prevFirstRangeIP.To4() != nil {
			bitLen = ipv4BitLen
		} else {
			bitLen = ipv6BitLen
		}
		_, _, right := partitionCIDR(spanningCIDR, net.IPNet{IP: prevFirstRangeIP, Mask: net.CIDRMask(bitLen, bitLen)})

		// Append all CIDRs but the first, as this CIDR includes the upper
		// bound of the spanning CIDR, which we still need to partition on.
		cidrList = append(cidrList, right...)
		spanningCIDR = *right[0]
		cidrList = cidrList[1:]
	}

	// Conversely, if the last IP of the spanning CIDR passes the upper bound
	// (lastIP), we need to split the spanning CIDR and only take the IPs that
	// are greater than the value which we split on, as we do not want the greater
	// values since they are greater than the upper-bound (lastIP).
	if bytes.Compare(*lastIPSpanning, lastIP) > 0 {
		// Split on the next IP of the last IP so that the left list of IPs
		// of the partition include the lastIP.
		nextFirstRangeIP := GetNextIP(lastIP)
		var bitLen int
		if nextFirstRangeIP.To4() != nil {
			bitLen = ipv4BitLen
		} else {
			bitLen = ipv6BitLen
		}
		left, _, _ := partitionCIDR(spanningCIDR, net.IPNet{IP: nextFirstRangeIP, Mask: net.CIDRMask(bitLen, bitLen)})
		cidrList = append(cidrList, left...)
	} else {
		// Otherwise, there is no need to partition; just use add the spanning
		// CIDR to the list of networks.
		cidrList = append(cidrList, &spanningCIDR)
	}
	return cidrList
}

// partitionCIDR returns a list of IP Networks partitioned upon excludeCIDR.
// The first list contains the networks to the left of the excludeCIDR in the
// partition,  the second is a list containing the excludeCIDR itself if it is
// contained within the targetCIDR (nil otherwise), and the
// third is a list containing the networks to the right of the excludeCIDR in
// the partition.
func partitionCIDR(targetCIDR, excludeCIDR net.IPNet) (left, excludeList, right []*net.IPNet) { //nolint
	var targetIsIPv4 bool
	if targetCIDR.IP.To4() != nil {
		targetIsIPv4 = true
	}

	targetIPRange := ipNetToRange(targetCIDR)
	excludeIPRange := ipNetToRange(excludeCIDR)

	targetFirstIP := *targetIPRange.First
	targetLastIP := *targetIPRange.Last

	excludeFirstIP := *excludeIPRange.First
	excludeLastIP := *excludeIPRange.Last

	targetMaskSize, _ := targetCIDR.Mask.Size()
	excludeMaskSize, _ := excludeCIDR.Mask.Size()

	if bytes.Compare(excludeLastIP, targetFirstIP) < 0 {
		return nil, nil, []*net.IPNet{&targetCIDR}
	} else if bytes.Compare(targetLastIP, excludeFirstIP) < 0 {
		return []*net.IPNet{&targetCIDR}, nil, nil
	}

	if targetMaskSize >= excludeMaskSize {
		return nil, []*net.IPNet{&targetCIDR}, nil
	}

	left = []*net.IPNet{}
	right = []*net.IPNet{}

	newPrefixLen := targetMaskSize + 1

	targetFirstCopy := make(net.IP, len(targetFirstIP))
	copy(targetFirstCopy, targetFirstIP)

	iLowerOld := make(net.IP, len(targetFirstCopy))
	copy(iLowerOld, targetFirstCopy)

	// Since golang only supports up to unsigned 64-bit integers, and we need
	// to perform addition on addresses, use math/big library, which allows
	// for manipulation of large integers.

	// Used to track the current lower and upper bounds of the ranges to compare
	// to excludeCIDR.
	iLower := big.NewInt(0)
	iUpper := big.NewInt(0)
	iLower = iLower.SetBytes(targetFirstCopy)

	var bitLen int

	if targetIsIPv4 {
		bitLen = ipv4BitLen
	} else {
		bitLen = ipv6BitLen
	}
	shiftAmount := uint(bitLen - newPrefixLen)

	targetIPInt := big.NewInt(0)
	targetIPInt.SetBytes(targetFirstIP.To16())

	exp := big.NewInt(0)

	// Use left shift for exponentiation
	exp = exp.Lsh(big.NewInt(1), shiftAmount)
	iUpper = iUpper.Add(targetIPInt, exp)

	matched := big.NewInt(0)

	for excludeMaskSize >= newPrefixLen {
		// Append leading zeros to IPv4 addresses, as math.Big.Int does not
		// append them when the IP address is copied from a byte array to
		// math.Big.Int. Leading zeroes are required for parsing IPv4 addresses
		// for use with net.IP / net.IPNet.
		var iUpperBytes, iLowerBytes []byte
		if targetIsIPv4 {
			iUpperBytes = append(ipv4LeadingZeroes, iUpper.Bytes()...) //nolint
			iLowerBytes = append(ipv4LeadingZeroes, iLower.Bytes()...) //nolint
		} else {
			iUpperBytesLen := len(iUpper.Bytes())
			// Make sure that the number of bytes in the array matches what net
			// package expects, as big package doesn't append leading zeroes.
			if iUpperBytesLen != net.IPv6len {
				numZeroesToAppend := net.IPv6len - iUpperBytesLen
				zeroBytes := make([]byte, numZeroesToAppend)
				iUpperBytes = append(zeroBytes, iUpper.Bytes()...) //nolint
			} else {
				iUpperBytes = iUpper.Bytes()
			}

			iLowerBytesLen := len(iLower.Bytes())
			if iLowerBytesLen != net.IPv6len {
				numZeroesToAppend := net.IPv6len - iLowerBytesLen
				zeroBytes := make([]byte, numZeroesToAppend)
				iLowerBytes = append(zeroBytes, iLower.Bytes()...) //nolint
			} else {
				iLowerBytes = iLower.Bytes()
			}
		}
		// If the IP we are excluding over is of a higher value than the current
		// CIDR prefix we are generating, add the CIDR prefix to the set of IPs
		// to the left of the exclude CIDR
		if bytes.Compare(excludeFirstIP, iUpperBytes) >= 0 {
			left = append(left, &net.IPNet{IP: iLowerBytes, Mask: net.CIDRMask(newPrefixLen, bitLen)})
			matched = matched.Set(iUpper)
		} else {
			// Same as above, but opposite.
			right = append(right, &net.IPNet{IP: iUpperBytes, Mask: net.CIDRMask(newPrefixLen, bitLen)})
			matched = matched.Set(iLower)
		}

		newPrefixLen++

		if newPrefixLen > bitLen {
			break
		}

		iLower = iLower.Set(matched)
		iUpper = iUpper.Add(matched, big.NewInt(0).Lsh(big.NewInt(1), uint(bitLen-newPrefixLen)))
	}
	excludeList = []*net.IPNet{&excludeCIDR}

	return left, excludeList, right
}

func RipeToCIDR(result string) (r []*net.IPNet) {
	var start, end net.IP

	ip_range := strings.ReplaceAll(result, " ", "")
	if idx := strings.IndexByte(ip_range, '-'); idx != -1 {
		start, end = net.ParseIP(ip_range[:idx]), net.ParseIP(ip_range[idx+1:])
	}

	return rangeToCIDRs(start, end)
}
