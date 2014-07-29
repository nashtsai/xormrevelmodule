Revel XORM Module
=================

This is Revel Module to enable [XORM](http://xorm.io) in the [Revel framework](http://revel.github.io/).

## Activation
To activate it, add the module to your app configuration:

	module.xorm = github.com/nashtsai/xormrevelmodule

## Options
This module takes DB modules options:

	db.import = github.com/go-sql-driver/mysql	# golang db driver
	db.driver = mysql							# driver name
	db.spec = root:@/mydb?charset=utf8			# datasource name

In addition you can set the maximum number of connections in the idle connection pool and maximum number of open connections to the database:

	db.maxidleconns = 10
	db.maxopenconns = 50

XORM specific options:

	xorm.showsql = true 	# show SQL
	xorm.showdebug = true	# show XORM debug info

## Using XORM controller
Instead of having anonymous *revel.Controller member replace with with *xormmodule.XormController

<pre class="prettyprint lang-go">
import (
	...
	"github.com/nashtsai/xormrevelmodule"
	"github.com/revel/revel"
	...
)

type MyXormController struct {
	*xormmodule.XormController
}

func (c MyXormController) List() revel.Result {
    users := make([]*Userinfo, 0)
    c.engine.Find(&users)
	return c.Render(users)
}
</pre>

## Post XORM engine init. handler

<pre class="prettyprint lang-go">
import (
	"github.com/go-xorm/xorm"
	"github.com/nashtsai/xormrevelmodule"
)

func init() {
	xormmodule.AddPostInitProcessor(func(engine *xorm.Engine){
		// your own init code, i.e., engine.Sync
    })
}

</pre>