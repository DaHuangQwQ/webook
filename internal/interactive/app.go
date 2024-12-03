package interactive

import (
	"github.com/DaHuangQwQ/gpkg/ginx"
	"github.com/DaHuangQwQ/gpkg/grpcx"
	"github.com/DaHuangQwQ/gpkg/saramax"
)

type App struct {
	Server    *grpcx.Server
	consumers []saramax.Consumer
	webAdmin  *ginx.Server
}
