package presenter

type Meta interface {
	GetHTTPStatus() int
}

// type Presenter interface {
// 	BuildSuccess(c *fiber.Ctx, data interface{}, message string, code int) error
// 	BuildError(c *fiber.Ctx, err error) error
// }
