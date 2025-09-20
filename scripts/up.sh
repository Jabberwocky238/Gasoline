iptables -A FORWARD -i wg0 -j ACCEPT; 
iptables -A FORWARD -o wg0 -j ACCEPT;
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE;

