package xormmodule

import (
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
)

type PostInitProcessorFunc func(*xorm.Engine)

var (
	Engine         *xorm.Engine
	Driver         string
	Spec           string
	MaxIdleConns   int
	MaxOpenConns   int
	ShowSQL        bool
	ShowDebug      bool
	postProcessors []PostInitProcessorFunc
)

func Init() {
	// Read configuration.
	var found bool
	if Driver, found = revel.Config.String("db.driver"); !found {
		revel.ERROR.Fatal("No db.driver found.")
	}
	if Spec, found = revel.Config.String("db.spec"); !found {
		revel.ERROR.Fatal("No db.spec found.")
	}
	// create Xorm engine
	var err error
	Engine, err = xorm.NewEngine(Driver, Spec)
	if err != nil {
		revel.ERROR.Fatal(err)
		return
	}

	if MaxIdleConns, found := revel.Config.Int("db.maxidleconns"); found {
		Engine.SetMaxIdleConns(MaxIdleConns)
	}

	if MaxOpenConns, found := revel.Config.Int("db.maxopenconns"); found {
		Engine.SetMaxOpenConns(MaxOpenConns)
	}

	if ShowSQL, found := revel.Config.Bool("xorm.showsql"); found {
		Engine.ShowSQL = ShowSQL
	}

	if ShowDebug, found := revel.Config.Bool("xorm.showdebug"); found {
		Engine.ShowDebug = ShowDebug
	}
	for _, processor := range postProcessors {
		processor(Engine)
	}
}

func AddPostInitProcessor(processor PostInitProcessorFunc) {
	if postProcessors == nil {
		postProcessors = []PostInitProcessorFunc{}
	}
	if processor != nil {
		postProcessors = append(postProcessors, processor)
	}
}

type XormController struct {
	*revel.Controller
	Engine *xorm.Engine
}

func (c *XormController) Attach() revel.Result {
	c.Engine = Engine
	return nil
}

// Rollback if it's still going (must have panicked).
func (c *XormController) Detach() revel.Result {
	c.Engine = nil
	return nil
}

func init() {
	revel.InterceptMethod((*XormController).Attach, revel.BEFORE)
	revel.InterceptMethod((*XormController).Detach, revel.FINALLY)

	revel.OnAppStart(Init)
}
