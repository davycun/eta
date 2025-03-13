package caller

import "github.com/davycun/eta/pkg/common/errs"

type Caller struct {
	stop bool
	Err  error
}

func NewCaller() *Caller {
	return &Caller{}
}
func (c *Caller) Call(f func(cl *Caller) error) *Caller {

	if c.Err != nil || c.stop {
		return c
	}
	c.Err = errs.Cover(c.Err, f(c))
	return c
}
func (c *Caller) Stop() *Caller {
	c.stop = true
	return c
}
