# spk aka spritzgebaeck

```
               _    
     ___ _ __ | | __
    / __| '_ \| |/ /
    \__ \ |_) |   < 
    |___/ .__/|_|\_\
         |_|         

        made with <3
```

*spk* aka **spritzgebaeck** is a small Golang tool that use different sources to get CIDRs/ASN for an organization.

# Building

* Download & install Go: https://golang.org/doc/install

```
go install github.com/dhn/spk@latest
```

# Usage

```
Usage of spk:
  -json
        Print results as JSON
  -s string
        Organization to find CIDRs for
  -silent
        Show only results in output
  -version
        Show version of spritzgebaeck
```

# Running spritzgebaeck

**spritzgebaeck** can be used to find CIDRs for a given organization. In this example for "TUI GmbH":

```
$ spk -silent -json -s "Dr. August Oetker KG"
{"cidr":"213.61.68.136/29","source":"ripe"}
{"cidr":"62.48.80.0/24","source":"ripe"}
{"cidr":"212.41.250.144/28","source":"ripe"}
{"cidr":"213.83.48.0/25","source":"ripe"}
{"cidr":"212.28.43.40/29","source":"ripe"}
{"cidr":"62.48.72.0/25","source":"ripe"}
{"cidr":"195.127.10.240/29","source":"ripe"}
{"cidr":"213.83.48.128/25","source":"ripe"}
{"cidr":"62.48.72.128/27","source":"ripe"}
{"cidr":"62.48.75.0/24","source":"ripe"}
{"cidr":"80.228.3.160/27","source":"ripe"}
{"cidr":"80.77.212.144/29","source":"ripe"}
{"cidr":"82.98.69.192/27","source":"ripe"}
{"cidr":"82.98.120.0/26","source":"ripe"}
{"cidr":"82.98.66.160/27","source":"ripe"}
{"cidr":"82.98.120.68/30","source":"ripe"}
{"cidr":"62.48.90.0/24","source":"ripe"}
{"cidr":"81.7.221.96/29","source":"ripe"}
{"cidr":"80.77.212.192/27","source":"ripe"}
{"cidr":"82.98.105.32/27","source":"ripe"}
{"cidr":"82.98.69.40/29","source":"ripe"}
{"cidr":"82.98.69.48/28","source":"ripe"}
{"cidr":"87.128.173.24/29","source":"ripe"}
{"cidr":"82.98.80.168/29","source":"ripe"}
{"cidr":"82.98.74.0/29","source":"ripe"}
{"cidr":"82.98.69.224/29","source":"ripe"}
{"cidr":"87.234.207.136/32","source":"ripe"}
{"cidr":"82.98.80.192/26","source":"ripe"}
{"cidr":"83.236.170.43/32","source":"ripe"}
{"cidr":"83.64.171.164/30","source":"ripe"}
{"cidr":"82.98.80.0/25","source":"ripe"}
{"cidr":"93.122.105.64/28","source":"ripe"}
{"cidr":"92.198.37.66/32","source":"ripe"}
{"cidr":"91.133.76.64/29","source":"ripe"}
{"cidr":"88.116.122.128/30","source":"ripe"}
{"cidr":"91.114.29.120/29","source":"ripe"}
{"cidr":"91.25.84.200/29","source":"ripe"}
```

## Some examples

Use `spk` to find a specific netrange + get more information through `whois`:

```
$ spk -silent -json -s "TUI GmbH" | xargs -I"{}" whois "{}"
% This is the RIPE Database query service.
% The objects are in RPSL format.
%
% The RIPE Database is subject to Terms and Conditions.
% See http://www.ripe.net/db/support/db-terms-conditions.pdf

% Note: this output has been filtered.
%       To receive output for a database update, use the "-B" flag.

% Information related to '213.61.68.136 - 213.61.68.143'

% Abuse contact for '213.61.68.136 - 213.61.68.143' is 'abuse@colt.net'

inetnum:        213.61.68.136 - 213.61.68.143
netname:        NET-DE-TUI-CRUISES-GMBH
descr:          TUI CRUISES GMBH
country:        DE
admin-c:        JO3608-RIPE
tech-c:         JO3608-RIPE
status:         ASSIGNED PA
mnt-by:         DE-COLT-MNT
created:        2015-06-05T11:18:06Z
last-modified:  2015-06-05T11:18:06Z
source:         RIPE

person:         JENS OETZEL
address:        TUI CRUISES GMBH
address:        BOUCHESTRASSE 12
address:        BERLIN, 12435, Germany
phone:          +49 40 600015380
nic-hdl:        JO3608-RIPE
mnt-by:         DE-COLT-MNT
created:        2015-06-05T11:18:06Z
last-modified:  2015-06-05T11:18:06Z
source:         RIPE

% Information related to '213.61.0.0/16AS8220'

route:          213.61.0.0/16
descr:          COLT-DE
origin:         AS8220
mnt-by:         DE-COLT-MNT
mnt-by:         MNT-COLT-SB
created:        2002-06-25T14:35:40Z
last-modified:  2022-07-12T15:01:06Z
source:         RIPE

% This query was served by the RIPE Database Query Service version 1.103 (WAGYU)


% This is the RIPE Database query service.
% The objects are in RPSL format.
%
% The RIPE Database is subject to Terms and Conditions.
% See http://www.ripe.net/db/support/db-terms-conditions.pdf

% Note: this output has been filtered.
%       To receive output for a database update, use the "-B" flag.

% Information related to '62.48.80.0 - 62.48.80.255'

% Abuse contact for '62.48.80.0 - 62.48.80.255' is 'abuse@net.de'

inetnum:        62.48.80.0 - 62.48.80.255
netname:        TUI-INFO-TEC
descr:          TUI InfoTec GmbH
descr:          Karl-Wiechert-Allee 4
descr:          30625 Hannover
country:        DE
admin-c:        CK2233-RIPE
tech-c:         CK2233-RIPE
status:         ASSIGNED PA
remarks:
mnt-by:         IPH-MNT
mnt-lower:      IPH-MNT
mnt-routes:     IPH-MNT
created:        2016-06-15T12:48:44Z
last-modified:  2016-06-15T12:48:44Z
source:         RIPE # Filtered

person:         Christian Kinzel
address:        TUI InfoTec GmbH
address:        Karl-Wiechert-Allee 4
address:        30625 Hannover
address:        Germany
phone:          +49 511 567-5148
fax-no:         +49 511 567-935148
nic-hdl:        CK2233-RIPE
mnt-by:         IPH-MNT
created:        2011-10-12T15:11:35Z
last-modified:  2015-03-09T10:11:56Z
source:         RIPE # Filtered

% Information related to '62.48.64.0/19AS15743'

route:          62.48.64.0/19
descr:          DE-IPH-20000621
origin:         AS15743
mnt-by:         IPH-MNT
created:        1970-01-01T00:00:00Z
last-modified:  2001-09-22T09:33:44Z
source:         RIPE

[...]
```

Let's use `spk` and `dnsx` [1] to find domains within the netranges:

```
$ spk -silent -s "TUI GmbH" | dnsx -silent -ptr -resp-only | grep tui
tuicrew-zqw-gw.zweibruecken.transkom.net
tuicrew-zqw-vpn.zweibruecken.transkom.net
kibana.tui-deutschland.plusline.de
itest-dhub-tui-deutschland.plusline.de
ctest.tui-uk.plusline.net
itest.tui-deutschland.plusline.de
ctest.tui-deutschland.plusline.de
ctest-dhub.tui-deutschland.plusline.de
webmail.tui.de
mail1.tui.de
tui-genius.tui.com
expiclub-bonussystem.tui.de
karatweb.tui.de
verkauf.tui.de
maritz.tui.de
mail.tui.de
auth.tui.de
webmail.tui.de
reservation-relay.tui.de
tlt.tui.de
products.tui.de
mail2.tui.de
abis.tui.com
www.tui.de
lanzarote.tui.com
ftp.tui.de
www.telis.tui.de
www.tui-personal.com
win1rob1.tui.de
www.tui-family.de
admin.tui.de
iriscms.tui.de
webtas.tui.com
tracetest.tui.com
www.tui-airlines.com
infotec.tui.de
mail.tuigroup.com
mailrelay.tuies.net
irisplus-internet.tui.de
mail.tui.de
one.tui.de
srv229.tui.de
tufis.tui.de
fisp.tui.de
meinschiff1.tuicruises.com
iasd.tui.de
meinschiff2.tuicruises.com
meinschiff.tuicruises.com
tui-blue.com
b-slave.tuicruises.plusline.de
dhub.tui-deutschland.plusline.de
prod.tui-deutschland.plusline.de
prod.tui-uk.plusline.net
```

Another use case could be to concat `dnsx`, `subfinder` [2] and `unfurl` [3] to find subdomains within the netrange:

```
$ spk -silent -s "TUI GmbH" | dnsx -silent -ptr -resp-only \
     | unfurl format %r.%t | sort -u | grep -i tui \
     | subfinder -silent -json -o output_subdomains.json

{"host":"www.tui-airlines.com","input":"tui-airlines.com","source":"riddler"}
{"host":"mail1.tui-airlines.com","input":"tui-airlines.com","source":"hackertarget"}
{"host":"timsconverter-stg-aws.es.tui.com","input":"tui.com","source":"alienvault"}
{"host":"webhookhub-invoicing.es.tui.com","input":"tui.com","source":"alienvault"}
{"host":"booking-in.lte.tui.com","input":"tui.com","source":"crtsh"}
{"host":"de-l13-eepool04-ws.tui.com","input":"tui.com","source":"crtsh"}
[...]
```

## References

- [1] https://github.com/projectdiscovery/subfinder
- [2] https://github.com/projectdiscovery/dnsx
- [3] https://github.com/tomnomnom/unfurl

### Trivia

[Spritzgebäck](https://en.wikipedia.org/wiki/Spritzgeb%C3%A4ck) means spritz cookies in German.