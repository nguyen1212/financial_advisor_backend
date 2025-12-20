package payload

import (
	"net/http"

	"github.com/gin-gonic/gin/binding"
)

func ShouldBindWith(
	req *http.Request,
	obj any,
	b binding.Binding,
) error {
	return b.Bind(req, obj)
}
