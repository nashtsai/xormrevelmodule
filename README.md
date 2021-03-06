Revel XORM Module
=================

This is Revel Module to enable [XORM](http://xorm.io) in the [Revel framework](http://revel.github.io/).

## Activation
To activate it, add the module to your app.conf:

	module.xorm = github.com/nashtsai/xormrevelmodule

## Options
This module takes Revel db module options:

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
Add anonymous xormmodule.XormController or anonymous xormmodule.XormSessionController member to your revel controller struct:

<pre class="prettyprint lang-go">
import (
	...
	"github.com/nashtsai/xormrevelmodule"
	"github.com/revel/revel"
	...
)

type MyXormController struct {
	*revel.Controller
	xormmodule.XormController
}

func (c MyXormController) List() revel.Result {
    users := make([]*Userinfo, 0)
    c.Engine.Find(&users)
	return c.Render(users)
}

type MyXormSessionController struct {
	*revel.Controller
	xormmodule.XormSessionController
}

func (c MyXormSessionController) Delete(id int64) revel.Result {
    _, err: = c.XormSession.Delete(&UserInfo{id:id})
	return c.Render(err)
}
</pre>

## Post XORM engine init. handler
Post init. handler after xorm.Engine is initialized:

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

## Advanced topics
[Go Walker API references](https://gowalker.org/github.com/nashtsai/xormrevelmodule)
