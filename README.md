# hlt 
Highlighter

Pipe data into it or "tail" a file by passing it as the final argument and it will highlight anything that matches any of the regular expressions specified with the -m option.

### Examples

tail -f /var/log/syslog | hlt -m http -c yellow -b

hlt -m http -c green -b -l /tmp/hlt.log /var/log/syslog
