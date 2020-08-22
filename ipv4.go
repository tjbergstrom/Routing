// ipv4.go


package router


import (
	"fmt"
	"net"
)


type IPv4 [4]byte


func ipv4_from_str(address string) (*IPv4, error) {
	addr := net.ParseIP(address)
	if addr == nil || addr.To4() == nil {
		return nil, fmt.Errorf("Address %s is not an IPv4 address", address)
	}
	ipv4 := IPv4{}
	for i:=0; i<4; i++ {
		ipv4[i] = addr[i]
	}
	return &ipv4, nil
}


// Is the address is a valid subnet mask
func (i *IPv4) is_mask() bool {
	ones := true
	var j uint
	for _, v := range i {
		if ones {
			// A whole byte of ones
			if v == 255 {
				continue
			}
			ones = false
			// Validate byte as having only leading ones
			for j = 0; j < 8; j++ {
				shifted := (v << j) & 255
				bit_is_zero := 128 != (shifted & 128)
				shifted_zero := 0 == shifted
				if !bit_is_zero {
					continue
				}
				if shifted_zero {
					break
				} else {
					// Found a zero bit before zeroing the byte
					return false
				}
			}
		} else if v != 0 {
			// Already had a 0 bit, so this byte should be all 0's
			return false
		}
	}
	return true
}


func (i *IPv4) mask_with_prefix(prefix uint8) *IPv4 {
	mask, err := mask_from_prefix(prefix)
	if err != nil {
		return nil
	}
	return i.mask_with(*mask)
}


func (i *IPv4) mask_with(mask IPv4) *IPv4 {
	if !mask.is_mask() {
		return nil
	}
	masked := &IPv4{}
	for j:=0; j<4; j++ {
		masked[j] = i[j] & mask[j]
	}
	return masked
}

func mask_from_prefix(prefix uint8) (*IPv4, error) {
	if prefix > 32 {
		return nil, fmt.Errorf("Prefix %d is greater than 32", prefix)
	}
	mask := &IPv4{}
	num_full := prefix / 8
	var i uint8
	for i = 0; i < num_full; i++ {
		mask[i] = 255
	}
	// Set value of last non-zero byte
	if num_full < 4 {
		rem := prefix % 8
		// 256 - (2 ^ (8-rem))
		power := 8 - rem
		partial := 1
		for i = 0; i < power; i++ {
			partial *= 2
		}
		mask[num_full] = byte(256 - partial)
	}
	return mask, nil
}



//
