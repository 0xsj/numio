// internal/eval/functions_network.go

package eval

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// IP PARSING & CONVERSION
// ════════════════════════════════════════════════════════════════

// FnIP converts an IP address string to its decimal representation.
// Args: IP address string (e.g., "192.168.1.1")
func FnIP(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("ip requires exactly 1 argument: IP address string")
	}

	ipStr := args[0].AsString()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return types.Errorf("ip: invalid IP address '%s'", ipStr)
	}

	// Convert to IPv4 if possible
	ipv4 := ip.To4()
	if ipv4 != nil {
		decimal := uint32(ipv4[0])<<24 | uint32(ipv4[1])<<16 | uint32(ipv4[2])<<8 | uint32(ipv4[3])
		return types.Number(float64(decimal))
	}

	// IPv6 - return as string representation of the full address
	return types.StringValue(ip.String())
}

// FnToIP converts a decimal number to an IP address string.
// Args: decimal number
func FnToIP(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("toip requires exactly 1 argument: decimal number")
	}

	decimal := uint32(args[0].AsFloat())

	ip := net.IPv4(
		byte(decimal>>24),
		byte(decimal>>16),
		byte(decimal>>8),
		byte(decimal),
	)

	return types.StringValue(ip.String())
}

// FnParseIP parses and normalizes an IP address.
// Args: IP address string
func FnParseIP(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("parseip requires exactly 1 argument: IP address string")
	}

	ipStr := args[0].AsString()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return types.Errorf("parseip: invalid IP address '%s'", ipStr)
	}

	return types.StringValue(ip.String())
}

// FnIPVersion returns the IP version (4 or 6).
func FnIPVersion(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("ipversion requires exactly 1 argument: IP address")
	}

	ipStr := args[0].AsString()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return types.Errorf("ipversion: invalid IP address '%s'", ipStr)
	}

	if ip.To4() != nil {
		return types.Number(4)
	}
	return types.Number(6)
}

// ════════════════════════════════════════════════════════════════
// IP CLASSIFICATION
// ════════════════════════════════════════════════════════════════

// FnIsPrivate checks if an IP address is in a private range.
func FnIsPrivate(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("isprivate requires exactly 1 argument: IP address")
	}

	ipStr := args[0].AsString()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return types.Errorf("isprivate: invalid IP address '%s'", ipStr)
	}

	if ip.IsPrivate() {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIsPublic checks if an IP address is public (not private, loopback, etc.).
func FnIsPublic(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("ispublic requires exactly 1 argument: IP address")
	}

	ipStr := args[0].AsString()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return types.Errorf("ispublic: invalid IP address '%s'", ipStr)
	}

	if !ip.IsPrivate() && !ip.IsLoopback() && !ip.IsMulticast() && !ip.IsUnspecified() && !ip.IsLinkLocalUnicast() {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIsLoopback checks if an IP address is a loopback address.
func FnIsLoopback(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("isloopback requires exactly 1 argument: IP address")
	}

	ipStr := args[0].AsString()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return types.Errorf("isloopback: invalid IP address '%s'", ipStr)
	}

	if ip.IsLoopback() {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIsMulticast checks if an IP address is a multicast address.
func FnIsMulticast(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("ismulticast requires exactly 1 argument: IP address")
	}

	ipStr := args[0].AsString()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return types.Errorf("ismulticast: invalid IP address '%s'", ipStr)
	}

	if ip.IsMulticast() {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIsLinkLocal checks if an IP address is link-local.
func FnIsLinkLocal(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("islinklocal requires exactly 1 argument: IP address")
	}

	ipStr := args[0].AsString()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return types.Errorf("islinklocal: invalid IP address '%s'", ipStr)
	}

	if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIPClass returns the class of an IPv4 address (A, B, C, D, E).
func FnIPClass(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("ipclass requires exactly 1 argument: IP address")
	}

	ipStr := args[0].AsString()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return types.Errorf("ipclass: invalid IP address '%s'", ipStr)
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return types.StringValue("IPv6")
	}

	firstOctet := ipv4[0]
	switch {
	case firstOctet < 128:
		return types.StringValue("A")
	case firstOctet < 192:
		return types.StringValue("B")
	case firstOctet < 224:
		return types.StringValue("C")
	case firstOctet < 240:
		return types.StringValue("D")
	default:
		return types.StringValue("E")
	}
}

// ════════════════════════════════════════════════════════════════
// SUBNET & CIDR FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnSubnet parses a CIDR and returns subnet information.
// Args: CIDR notation (e.g., "192.168.1.0/24")
func FnSubnet(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("subnet requires exactly 1 argument: CIDR notation")
	}

	cidrStr := args[0].AsString()
	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return types.Errorf("subnet: invalid CIDR '%s'", cidrStr)
	}

	ones, bits := ipNet.Mask.Size()
	hostBits := bits - ones
	var hosts uint64
	if hostBits >= 64 {
		hosts = 0 // Too large to represent
	} else if hostBits <= 1 {
		hosts = 0 // /31 or /32
	} else {
		hosts = (1 << hostBits) - 2 // Subtract network and broadcast
	}

	return types.Number(float64(hosts))
}

// FnHosts returns the number of usable hosts in a subnet.
// Args: CIDR notation or prefix length
func FnHosts(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("hosts requires exactly 1 argument: CIDR or prefix length")
	}

	// Check if it's a number (prefix length) or string (CIDR)
	if args[0].IsNumber() {
		prefix := int(args[0].AsFloat())
		if prefix < 0 || prefix > 32 {
			return types.Error("hosts: prefix must be between 0 and 32")
		}
		hostBits := 32 - prefix
		if hostBits <= 1 {
			return types.Number(0)
		}
		hosts := (1 << hostBits) - 2
		return types.Number(float64(hosts))
	}

	// Parse as CIDR
	cidrStr := args[0].AsString()
	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return types.Errorf("hosts: invalid CIDR '%s'", cidrStr)
	}

	ones, bits := ipNet.Mask.Size()
	hostBits := bits - ones
	if hostBits <= 1 {
		return types.Number(0)
	}
	hosts := (1 << hostBits) - 2

	return types.Number(float64(hosts))
}

// FnNetmask converts a prefix length to a netmask or vice versa.
// Args: prefix length (number) or netmask (string)
func FnNetmask(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("netmask requires exactly 1 argument: prefix length or netmask")
	}

	// If number, convert prefix to netmask
	if args[0].IsNumber() {
		prefix := int(args[0].AsFloat())
		if prefix < 0 || prefix > 32 {
			return types.Error("netmask: prefix must be between 0 and 32")
		}

		mask := net.CIDRMask(prefix, 32)
		return types.StringValue(net.IP(mask).String())
	}

	// If string, convert netmask to prefix
	maskStr := args[0].AsString()

	// Handle /24 notation
	if strings.HasPrefix(maskStr, "/") {
		prefix, err := strconv.Atoi(maskStr[1:])
		if err != nil || prefix < 0 || prefix > 32 {
			return types.Errorf("netmask: invalid prefix '%s'", maskStr)
		}
		mask := net.CIDRMask(prefix, 32)
		return types.StringValue(net.IP(mask).String())
	}

	// Parse as IP-style netmask
	ip := net.ParseIP(maskStr)
	if ip == nil {
		return types.Errorf("netmask: invalid netmask '%s'", maskStr)
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return types.Error("netmask: only IPv4 netmasks supported")
	}

	mask := net.IPMask(ipv4)
	ones, _ := mask.Size()

	return types.Number(float64(ones))
}

// FnCIDR creates a CIDR notation from IP and netmask.
// Args: IP address, netmask or prefix
func FnCIDR(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("cidr requires 2 arguments: IP address, netmask/prefix")
	}

	ipStr := args[0].AsString()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return types.Errorf("cidr: invalid IP address '%s'", ipStr)
	}

	var prefix int

	if args[1].IsNumber() {
		prefix = int(args[1].AsFloat())
	} else {
		maskStr := args[1].AsString()

		// Handle /24 notation
		if strings.HasPrefix(maskStr, "/") {
			p, err := strconv.Atoi(maskStr[1:])
			if err != nil {
				return types.Errorf("cidr: invalid prefix '%s'", maskStr)
			}
			prefix = p
		} else {
			// Parse as netmask
			maskIP := net.ParseIP(maskStr)
			if maskIP == nil {
				return types.Errorf("cidr: invalid netmask '%s'", maskStr)
			}
			ipv4 := maskIP.To4()
			if ipv4 == nil {
				return types.Error("cidr: only IPv4 supported")
			}
			mask := net.IPMask(ipv4)
			prefix, _ = mask.Size()
		}
	}

	if prefix < 0 || prefix > 32 {
		return types.Error("cidr: prefix must be between 0 and 32")
	}

	return types.StringValue(fmt.Sprintf("%s/%d", ip.String(), prefix))
}

// FnNetwork returns the network address for a CIDR.
func FnNetwork(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("network requires exactly 1 argument: CIDR notation")
	}

	cidrStr := args[0].AsString()
	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return types.Errorf("network: invalid CIDR '%s'", cidrStr)
	}

	return types.StringValue(ipNet.IP.String())
}

// FnBroadcast returns the broadcast address for a CIDR.
func FnBroadcast(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("broadcast requires exactly 1 argument: CIDR notation")
	}

	cidrStr := args[0].AsString()
	ip, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return types.Errorf("broadcast: invalid CIDR '%s'", cidrStr)
	}

	// Calculate broadcast: network OR (NOT mask)
	ipv4 := ip.To4()
	if ipv4 == nil {
		return types.Error("broadcast: only IPv4 supported")
	}

	mask := ipNet.Mask
	broadcast := make(net.IP, len(ipv4))
	for i := range ipv4 {
		broadcast[i] = ipNet.IP[i] | ^mask[i]
	}

	return types.StringValue(broadcast.String())
}

// FnFirstHost returns the first usable host address in a subnet.
func FnFirstHost(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("firsthost requires exactly 1 argument: CIDR notation")
	}

	cidrStr := args[0].AsString()
	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return types.Errorf("firsthost: invalid CIDR '%s'", cidrStr)
	}

	ipv4 := ipNet.IP.To4()
	if ipv4 == nil {
		return types.Error("firsthost: only IPv4 supported")
	}

	// First host is network + 1
	firstHost := make(net.IP, len(ipv4))
	copy(firstHost, ipv4)
	firstHost[3]++

	return types.StringValue(firstHost.String())
}

// FnLastHost returns the last usable host address in a subnet.
func FnLastHost(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("lasthost requires exactly 1 argument: CIDR notation")
	}

	cidrStr := args[0].AsString()
	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return types.Errorf("lasthost: invalid CIDR '%s'", cidrStr)
	}

	ipv4 := ipNet.IP.To4()
	if ipv4 == nil {
		return types.Error("lasthost: only IPv4 supported")
	}

	mask := ipNet.Mask

	// Last host is broadcast - 1
	lastHost := make(net.IP, len(ipv4))
	for i := range ipv4 {
		lastHost[i] = ipv4[i] | ^mask[i]
	}
	lastHost[3]--

	return types.StringValue(lastHost.String())
}

// FnWildcard returns the wildcard mask (inverse of netmask).
func FnWildcard(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("wildcard requires exactly 1 argument: prefix length or netmask")
	}

	var mask net.IPMask

	if args[0].IsNumber() {
		prefix := int(args[0].AsFloat())
		if prefix < 0 || prefix > 32 {
			return types.Error("wildcard: prefix must be between 0 and 32")
		}
		mask = net.CIDRMask(prefix, 32)
	} else {
		maskStr := args[0].AsString()

		if strings.HasPrefix(maskStr, "/") {
			prefix, err := strconv.Atoi(maskStr[1:])
			if err != nil || prefix < 0 || prefix > 32 {
				return types.Errorf("wildcard: invalid prefix '%s'", maskStr)
			}
			mask = net.CIDRMask(prefix, 32)
		} else {
			ip := net.ParseIP(maskStr)
			if ip == nil {
				return types.Errorf("wildcard: invalid netmask '%s'", maskStr)
			}
			ipv4 := ip.To4()
			if ipv4 == nil {
				return types.Error("wildcard: only IPv4 supported")
			}
			mask = net.IPMask(ipv4)
		}
	}

	// Wildcard is inverse of mask
	wildcard := make(net.IP, len(mask))
	for i := range mask {
		wildcard[i] = ^mask[i]
	}

	return types.StringValue(wildcard.String())
}

// ════════════════════════════════════════════════════════════════
// IP MATH & RANGE
// ════════════════════════════════════════════════════════════════

// FnIPAdd adds an offset to an IP address.
// Args: IP address, offset
func FnIPAdd(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("ipadd requires 2 arguments: IP address, offset")
	}

	ipStr := args[0].AsString()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return types.Errorf("ipadd: invalid IP address '%s'", ipStr)
	}

	offset := int64(args[1].AsFloat())

	ipv4 := ip.To4()
	if ipv4 == nil {
		return types.Error("ipadd: only IPv4 supported")
	}

	decimal := int64(uint32(ipv4[0])<<24 | uint32(ipv4[1])<<16 | uint32(ipv4[2])<<8 | uint32(ipv4[3]))
	decimal += offset

	if decimal < 0 || decimal > 0xFFFFFFFF {
		return types.Error("ipadd: result out of range")
	}

	result := net.IPv4(
		byte(decimal>>24),
		byte(decimal>>16),
		byte(decimal>>8),
		byte(decimal),
	)

	return types.StringValue(result.String())
}

// FnIPSub subtracts an offset from an IP address or calculates difference.
// Args: IP address, offset OR two IP addresses
func FnIPSub(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("ipsub requires 2 arguments: IP1, IP2 or offset")
	}

	ip1Str := args[0].AsString()
	ip1 := net.ParseIP(ip1Str)
	if ip1 == nil {
		return types.Errorf("ipsub: invalid IP address '%s'", ip1Str)
	}

	ipv4_1 := ip1.To4()
	if ipv4_1 == nil {
		return types.Error("ipsub: only IPv4 supported")
	}

	decimal1 := int64(uint32(ipv4_1[0])<<24 | uint32(ipv4_1[1])<<16 | uint32(ipv4_1[2])<<8 | uint32(ipv4_1[3]))

	// Check if second arg is IP or number
	if args[1].IsNumber() {
		offset := int64(args[1].AsFloat())
		decimal1 -= offset

		if decimal1 < 0 || decimal1 > 0xFFFFFFFF {
			return types.Error("ipsub: result out of range")
		}

		result := net.IPv4(
			byte(decimal1>>24),
			byte(decimal1>>16),
			byte(decimal1>>8),
			byte(decimal1),
		)
		return types.StringValue(result.String())
	}

	// Second arg is an IP - calculate difference
	ip2Str := args[1].AsString()
	ip2 := net.ParseIP(ip2Str)
	if ip2 == nil {
		return types.Errorf("ipsub: invalid IP address '%s'", ip2Str)
	}

	ipv4_2 := ip2.To4()
	if ipv4_2 == nil {
		return types.Error("ipsub: only IPv4 supported")
	}

	decimal2 := int64(uint32(ipv4_2[0])<<24 | uint32(ipv4_2[1])<<16 | uint32(ipv4_2[2])<<8 | uint32(ipv4_2[3]))

	return types.Number(float64(decimal1 - decimal2))
}

// FnIPInRange checks if an IP is within a CIDR range.
// Args: IP address, CIDR notation
func FnIPInRange(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("ipinrange requires 2 arguments: IP address, CIDR")
	}

	ipStr := args[0].AsString()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return types.Errorf("ipinrange: invalid IP address '%s'", ipStr)
	}

	cidrStr := args[1].AsString()
	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return types.Errorf("ipinrange: invalid CIDR '%s'", cidrStr)
	}

	if ipNet.Contains(ip) {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIPRange returns the number of addresses in a range.
// Args: start IP, end IP
func FnIPRange(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("iprange requires 2 arguments: start IP, end IP")
	}

	ip1Str := args[0].AsString()
	ip1 := net.ParseIP(ip1Str)
	if ip1 == nil {
		return types.Errorf("iprange: invalid IP address '%s'", ip1Str)
	}

	ip2Str := args[1].AsString()
	ip2 := net.ParseIP(ip2Str)
	if ip2 == nil {
		return types.Errorf("iprange: invalid IP address '%s'", ip2Str)
	}

	ipv4_1 := ip1.To4()
	ipv4_2 := ip2.To4()
	if ipv4_1 == nil || ipv4_2 == nil {
		return types.Error("iprange: only IPv4 supported")
	}

	decimal1 := uint32(ipv4_1[0])<<24 | uint32(ipv4_1[1])<<16 | uint32(ipv4_1[2])<<8 | uint32(ipv4_1[3])
	decimal2 := uint32(ipv4_2[0])<<24 | uint32(ipv4_2[1])<<16 | uint32(ipv4_2[2])<<8 | uint32(ipv4_2[3])

	var diff uint32
	if decimal2 >= decimal1 {
		diff = decimal2 - decimal1 + 1
	} else {
		diff = decimal1 - decimal2 + 1
	}

	return types.Number(float64(diff))
}

// ════════════════════════════════════════════════════════════════
// IP COMPONENT EXTRACTION
// ════════════════════════════════════════════════════════════════

// FnOctet extracts an octet from an IPv4 address.
// Args: IP address, octet index (1-4)
func FnOctet(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("octet requires 2 arguments: IP address, index (1-4)")
	}

	ipStr := args[0].AsString()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return types.Errorf("octet: invalid IP address '%s'", ipStr)
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return types.Error("octet: only IPv4 supported")
	}

	index := int(args[1].AsFloat())
	if index < 1 || index > 4 {
		return types.Error("octet: index must be 1-4")
	}

	return types.Number(float64(ipv4[index-1]))
}

// FnIPBinary returns the binary representation of an IP.
func FnIPBinary(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("ipbinary requires exactly 1 argument: IP address")
	}

	ipStr := args[0].AsString()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return types.Errorf("ipbinary: invalid IP address '%s'", ipStr)
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return types.Error("ipbinary: only IPv4 supported")
	}

	binary := fmt.Sprintf("%08b.%08b.%08b.%08b", ipv4[0], ipv4[1], ipv4[2], ipv4[3])
	return types.StringValue(binary)
}

// ════════════════════════════════════════════════════════════════
// COMMON NETWORK CALCULATIONS
// ════════════════════════════════════════════════════════════════

// FnSubnets calculates how many subnets can be created.
// Args: original prefix, new prefix
func FnSubnets(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("subnets requires 2 arguments: original prefix, new prefix")
	}

	original := int(args[0].AsFloat())
	new := int(args[1].AsFloat())

	if original < 0 || original > 32 || new < 0 || new > 32 {
		return types.Error("subnets: prefixes must be between 0 and 32")
	}

	if new < original {
		return types.Error("subnets: new prefix must be >= original prefix")
	}

	borrowed := new - original
	subnets := 1 << borrowed

	return types.Number(float64(subnets))
}

// FnPrefixFromHosts calculates the prefix length needed for a given number of hosts.
func FnPrefixFromHosts(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("prefixfromhosts requires exactly 1 argument: number of hosts")
	}

	hosts := int(args[0].AsFloat())
	if hosts < 1 {
		return types.Number(32)
	}

	// Need hosts + 2 addresses (network + broadcast)
	needed := hosts + 2

	// Find minimum host bits
	hostBits := 0
	for (1 << hostBits) < needed {
		hostBits++
	}

	prefix := 32 - hostBits
	if prefix < 0 {
		prefix = 0
	}

	return types.Number(float64(prefix))
}

// FnSameSubnet checks if two IPs are in the same subnet.
// Args: IP1, IP2, netmask or prefix
func FnSameSubnet(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("samesubnet requires 3 arguments: IP1, IP2, netmask/prefix")
	}

	ip1Str := args[0].AsString()
	ip1 := net.ParseIP(ip1Str)
	if ip1 == nil {
		return types.Errorf("samesubnet: invalid IP address '%s'", ip1Str)
	}

	ip2Str := args[1].AsString()
	ip2 := net.ParseIP(ip2Str)
	if ip2 == nil {
		return types.Errorf("samesubnet: invalid IP address '%s'", ip2Str)
	}

	var mask net.IPMask

	if args[2].IsNumber() {
		prefix := int(args[2].AsFloat())
		if prefix < 0 || prefix > 32 {
			return types.Error("samesubnet: prefix must be between 0 and 32")
		}
		mask = net.CIDRMask(prefix, 32)
	} else {
		maskStr := args[2].AsString()
		maskIP := net.ParseIP(maskStr)
		if maskIP == nil {
			return types.Errorf("samesubnet: invalid netmask '%s'", maskStr)
		}
		ipv4 := maskIP.To4()
		if ipv4 == nil {
			return types.Error("samesubnet: only IPv4 supported")
		}
		mask = net.IPMask(ipv4)
	}

	ipv4_1 := ip1.To4()
	ipv4_2 := ip2.To4()
	if ipv4_1 == nil || ipv4_2 == nil {
		return types.Error("samesubnet: only IPv4 supported")
	}

	// Compare network addresses
	for i := range ipv4_1 {
		if (ipv4_1[i] & mask[i]) != (ipv4_2[i] & mask[i]) {
			return types.Number(0)
		}
	}

	return types.Number(1)
}
