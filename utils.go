package goproxmox

type BitBool bool

func (b *BitBool) UnmarshalJSON(bytes []byte) error {
	*b = string(bytes) == "1"
	return nil
}
