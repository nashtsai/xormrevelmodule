package xormmodule

import (
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
)

type PostInitProcessorFunc func(*xorm.Engine)
type SessionHandlerFunc func(*xorm.Session) error

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

// XormController to be added as anonymous member to the revel controller struct, with XormController.Engine attached.
type XormController struct {
	Engine      *xorm.Engine
	XormSession *xorm.Session
}

// Create xorm.Session and attached to XormSession member if not already attached.
func (c *XormController) AttachSession() {
	if c.XormSession == nil {
		c.XormSession = c.Engine.NewSession()
	}
}

// Detach XormSession member and call xorm.Session.Close() if has attached XormSession.
func (c *XormController) DetachSession() {
	if c.XormSession != nil {
		c.XormSession.Close()
		c.XormSession = nil
	}
}

// Attach XormSession and call handler.
func (c *XormController) WithSession(handler SessionHandlerFunc) error {
	c.AttachSession()
	return handler(c.XormSession)
}

// Create a new xorm.Session and call handler, the xorm.Session will be closed after handler called.
func (c *XormController) WithNewSession(handler SessionHandlerFunc) error {
	session := c.Engine.NewSession()
	defer session.Close()
	return handler(session)
}

func (c *XormController) doTransaction(session *xorm.Session, handler SessionHandlerFunc) error {
	err := session.Begin()
	if err != nil {
		return err
	}

	if err = handler(session); err == nil {
		session.Commit()
	} else {
		session.Rollback()
	}
	return err
}

// Begin a SQL transaction and if handler did not return error
// it will commit the transaction, otherwise rollback the transaction.
// This will use attached XormSession if already attached.
func (c *XormController) WithTx(handler SessionHandlerFunc) error {
	if c.XormSession != nil {
		return c.doTransaction(c.XormSession, handler)
	} else {
		return c.WithNewTx(handler)
	}

}

// Begin a SQL transaction and if handler did not return error
// it will commit the transaction, otherwise rollback the transaction.
// This will create a new xorm.Session for the handler.
func (c *XormController) WithNewTx(handler SessionHandlerFunc) error {
	session := c.Engine.NewSession()
	defer session.Close()
	return c.doTransaction(session, handler)
}

// Attach XormController.Engine, this is automatically done at revel.BEFORE revel.InterceptMethod step.
func (c *XormController) Attach() revel.Result {
	revel.TRACE.Printf("(*XormController).Attach")
	c.Engine = Engine
	return nil
}

// Detach XormController.Engine, this is automatically done at revel.FINALLY revel.InterceptMethod step.
func (c *XormController) Detach() revel.Result {
	revel.TRACE.Printf("(*XormController).Detach")
	if c.XormSession != nil {
		c.XormSession.Rollback()
		c.XormSession.Close()
		c.XormSession = nil
	}
	c.Engine = nil
	return nil
}

// Commit and Close attached XormSession, this is automatically done at revel.AFTER revel.InterceptMethod step.
// Issue panic if commit XormSession is undesired.
func (c *XormController) Commit() revel.Result {
	revel.TRACE.Printf("(*XormController).Commit")
	if c.XormSession != nil {
		c.XormSession.Commit()
		c.XormSession.Close()
		c.XormSession = nil
	}
	return nil
}

// XormSessionController to be added as anonymous member to the revel controller struct, with XormController.Engine
// and XormController.XormSession attached.
type XormSessionController struct {
	XormController
}

// Attach XormSessionController.XormSession, this is automatically done at revel.BEFORE revel.InterceptMethod step.
func (c *XormSessionController) Attach() revel.Result {
	revel.TRACE.Printf("(*XormSessionController).Attach")
	c.XormSession = c.Engine.NewSession()
	return nil
}

func init() {
	revel.InterceptMethod((*XormController).Attach, revel.BEFORE)
	revel.InterceptMethod((*XormController).Commit, revel.AFTER)
	revel.InterceptMethod((*XormController).Detach, revel.FINALLY)

	revel.InterceptMethod((*XormSessionController).Attach, revel.BEFORE)

	revel.OnAppStart(Init)
}
