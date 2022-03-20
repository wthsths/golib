package gl_readonly

type String struct {
	value string
}

func NewString(value string) String {
	return String{
		value: value,
	}
}

func (v String) Value() string {
	return v.value
}
