package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
	"runtime"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var db *sql.DB
var (
	UserLockThreshold int
	IPBanThreshold    int
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func init() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		getEnv("ISU4_DB_USER", "root"),
		getEnv("ISU4_DB_PASSWORD", ""),
		getEnv("ISU4_DB_HOST", "localhost"),
		getEnv("ISU4_DB_PORT", "3306"),
		getEnv("ISU4_DB_NAME", "isu4_qualifier"),
	)

	var err error

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	UserLockThreshold, err = strconv.Atoi(getEnv("ISU4_USER_LOCK_THRESHOLD", "3"))
	if err != nil {
		panic(err)
	}

	IPBanThreshold, err = strconv.Atoi(getEnv("ISU4_IP_BAN_THRESHOLD", "10"))
	if err != nil {
		panic(err)
	}
}


func getIndex(c *gin.Context) {
	query := c.Request.URL.Query()
	param := query.Get("err")

	if param == "banned" {
		c.Data(200, "text/html", []byte(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <script id="css">
    document.write('<link rel="stylesheet" href="/stylesheets/bootstrap.min.css"> <link rel="stylesheet" href="/stylesheets/bootflat.min.css"> <link rel="stylesheet" href="/stylesheets/isucon-bank.css">');
    script = document.getElementById('css');
    script.parentNode.removeChild(script);
    </script>
    <title>isucon4</title>
  </head>
  <body>
    <div class="container">
      <h1 id="topbar">
        <a href="/">
          <script id="img">
          document.write('<img src="/images/isucon-bank.png" alt="いすこん銀行 オンラインバンキングサービス">');
          script = document.getElementById('img');
          script.parentNode.removeChild(script);
          </script>
        </a>
      </h1>
      <div id="be-careful-phising" class="panel panel-danger">
  <div class="panel-heading">
    <span class="hikaru-mozi">偽画面にご注意ください！</span>
  </div>
  <div class="panel-body">
    <p>偽のログイン画面を表示しお客様の情報を盗み取ろうとする犯罪が多発しています。</p>
    <p>ログイン直後にダウンロード中や、見知らぬウィンドウが開いた場合、<br>すでにウィルスに感染している場合がございます。即座に取引を中止してください。</p>
    <p>また、残高照会のみなど、必要のない場面で乱数表の入力を求められても、<br>絶対に入力しないでください。</p>
  </div>
</div>

<div class="page-header">
  <h1>ログイン</h1>
</div>

  <div id="notice-message" class="alert alert-danger" role="alert">You're banned.</div>

<div class="container">
  <form class="form-horizontal" role="form" action="/login" method="POST">
    <div class="form-group">
      <label for="input-username" class="col-sm-3 control-label">お客様ご契約ID</label>
      <div class="col-sm-9">
        <input id="input-username" type="text" class="form-control" placeholder="半角英数字" name="login">
      </div>
    </div>
    <div class="form-group">
      <label for="input-password" class="col-sm-3 control-label">パスワード</label>
      <div class="col-sm-9">
        <input type="password" class="form-control" id="input-password" name="password" placeholder="半角英数字・記号（２文字以上）">
      </div>
    </div>
    <div class="form-group">
      <div class="col-sm-offset-3 col-sm-9">
        <button type="submit" class="btn btn-primary btn-lg btn-block">ログイン</button>
      </div>
    </div>
  </form>
</div>

    </div>

  </body>
</html>
		`))
	} else if param == "wrong" {
		c.Data(200, "text/html", []byte(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <script id="css">
    document.write('<link rel="stylesheet" href="/stylesheets/bootstrap.min.css"> <link rel="stylesheet" href="/stylesheets/bootflat.min.css"> <link rel="stylesheet" href="/stylesheets/isucon-bank.css">');
    script = document.getElementById('css');
    script.parentNode.removeChild(script);
    </script>
    <title>isucon4</title>
  </head>
  <body>
    <div class="container">
      <h1 id="topbar">
        <a href="/">
          <script id="img">
          document.write('<img src="/images/isucon-bank.png" alt="いすこん銀行 オンラインバンキングサービス">');
          script = document.getElementById('img');
          script.parentNode.removeChild(script);
          </script>
        </a>
      </h1>
      <div id="be-careful-phising" class="panel panel-danger">
  <div class="panel-heading">
    <span class="hikaru-mozi">偽画面にご注意ください！</span>
  </div>
  <div class="panel-body">
    <p>偽のログイン画面を表示しお客様の情報を盗み取ろうとする犯罪が多発しています。</p>
    <p>ログイン直後にダウンロード中や、見知らぬウィンドウが開いた場合、<br>すでにウィルスに感染している場合がございます。即座に取引を中止してください。</p>
    <p>また、残高照会のみなど、必要のない場面で乱数表の入力を求められても、<br>絶対に入力しないでください。</p>
  </div>
</div>

<div class="page-header">
  <h1>ログイン</h1>
</div>

  <div id="notice-message" class="alert alert-danger" role="alert">Wrong username or password</div>

<div class="container">
  <form class="form-horizontal" role="form" action="/login" method="POST">
    <div class="form-group">
      <label for="input-username" class="col-sm-3 control-label">お客様ご契約ID</label>
      <div class="col-sm-9">
        <input id="input-username" type="text" class="form-control" placeholder="半角英数字" name="login">
      </div>
    </div>
    <div class="form-group">
      <label for="input-password" class="col-sm-3 control-label">パスワード</label>
      <div class="col-sm-9">
        <input type="password" class="form-control" id="input-password" name="password" placeholder="半角英数字・記号（２文字以上）">
      </div>
    </div>
    <div class="form-group">
      <div class="col-sm-offset-3 col-sm-9">
        <button type="submit" class="btn btn-primary btn-lg btn-block">ログイン</button>
      </div>
    </div>
  </form>
</div>

    </div>

  </body>
</html>
		`))
	} else if param == "invalid" {
		c.Data(200, "text/html", []byte(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <script id="css">
    document.write('<link rel="stylesheet" href="/stylesheets/bootstrap.min.css"> <link rel="stylesheet" href="/stylesheets/bootflat.min.css"> <link rel="stylesheet" href="/stylesheets/isucon-bank.css">');
    script = document.getElementById('css');
    script.parentNode.removeChild(script);
    </script>
    <title>isucon4</title>
  </head>
  <body>
    <div class="container">
      <h1 id="topbar">
        <a href="/">
          <script id="img">
          document.write('<img src="/images/isucon-bank.png" alt="いすこん銀行 オンラインバンキングサービス">');
          script = document.getElementById('img');
          script.parentNode.removeChild(script);
          </script>
        </a>
      </h1>
      <div id="be-careful-phising" class="panel panel-danger">
  <div class="panel-heading">
    <span class="hikaru-mozi">偽画面にご注意ください！</span>
  </div>
  <div class="panel-body">
    <p>偽のログイン画面を表示しお客様の情報を盗み取ろうとする犯罪が多発しています。</p>
    <p>ログイン直後にダウンロード中や、見知らぬウィンドウが開いた場合、<br>すでにウィルスに感染している場合がございます。即座に取引を中止してください。</p>
    <p>また、残高照会のみなど、必要のない場面で乱数表の入力を求められても、<br>絶対に入力しないでください。</p>
  </div>
</div>

<div class="page-header">
  <h1>ログイン</h1>
</div>

  <div id="notice-message" class="alert alert-danger" role="alert">You must be logged in</div>

<div class="container">
  <form class="form-horizontal" role="form" action="/login" method="POST">
    <div class="form-group">
      <label for="input-username" class="col-sm-3 control-label">お客様ご契約ID</label>
      <div class="col-sm-9">
        <input id="input-username" type="text" class="form-control" placeholder="半角英数字" name="login">
      </div>
    </div>
    <div class="form-group">
      <label for="input-password" class="col-sm-3 control-label">パスワード</label>
      <div class="col-sm-9">
        <input type="password" class="form-control" id="input-password" name="password" placeholder="半角英数字・記号（２文字以上）">
      </div>
    </div>
    <div class="form-group">
      <div class="col-sm-offset-3 col-sm-9">
        <button type="submit" class="btn btn-primary btn-lg btn-block">ログイン</button>
      </div>
    </div>
  </form>
</div>

    </div>

  </body>
</html>
		`))
	} else if param == "locked" {
		c.Data(200, "text/html", []byte(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <script id="css">
    document.write('<link rel="stylesheet" href="/stylesheets/bootstrap.min.css"> <link rel="stylesheet" href="/stylesheets/bootflat.min.css"> <link rel="stylesheet" href="/stylesheets/isucon-bank.css">');
    script = document.getElementById('css');
    script.parentNode.removeChild(script);
    </script>
    <title>isucon4</title>
  </head>
  <body>
    <div class="container">
      <h1 id="topbar">
        <a href="/">
          <script id="img">
          document.write('<img src="/images/isucon-bank.png" alt="いすこん銀行 オンラインバンキングサービス">');
          script = document.getElementById('img');
          script.parentNode.removeChild(script);
          </script>
        </a>
      </h1>
      <div id="be-careful-phising" class="panel panel-danger">
  <div class="panel-heading">
    <span class="hikaru-mozi">偽画面にご注意ください！</span>
  </div>
  <div class="panel-body">
    <p>偽のログイン画面を表示しお客様の情報を盗み取ろうとする犯罪が多発しています。</p>
    <p>ログイン直後にダウンロード中や、見知らぬウィンドウが開いた場合、<br>すでにウィルスに感染している場合がございます。即座に取引を中止してください。</p>
    <p>また、残高照会のみなど、必要のない場面で乱数表の入力を求められても、<br>絶対に入力しないでください。</p>
  </div>
</div>

<div class="page-header">
  <h1>ログイン</h1>
</div>

  <div id="notice-message" class="alert alert-danger" role="alert">This account is locked.</div>

<div class="container">
  <form class="form-horizontal" role="form" action="/login" method="POST">
    <div class="form-group">
      <label for="input-username" class="col-sm-3 control-label">お客様ご契約ID</label>
      <div class="col-sm-9">
        <input id="input-username" type="text" class="form-control" placeholder="半角英数字" name="login">
      </div>
    </div>
    <div class="form-group">
      <label for="input-password" class="col-sm-3 control-label">パスワード</label>
      <div class="col-sm-9">
        <input type="password" class="form-control" id="input-password" name="password" placeholder="半角英数字・記号（２文字以上）">
      </div>
    </div>
    <div class="form-group">
      <div class="col-sm-offset-3 col-sm-9">
        <button type="submit" class="btn btn-primary btn-lg btn-block">ログイン</button>
      </div>
    </div>
  </form>
</div>

    </div>

  </body>
</html>
		`))
	} else {
		c.Data(200, "text/html", []byte(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <script id="css">
    document.write('<link rel="stylesheet" href="/stylesheets/bootstrap.min.css"> <link rel="stylesheet" href="/stylesheets/bootflat.min.css"> <link rel="stylesheet" href="/stylesheets/isucon-bank.css">');
    script = document.getElementById('css');
    script.parentNode.removeChild(script);
    </script>
    <title>isucon4</title>
  </head>
  <body>
    <div class="container">
      <h1 id="topbar">
        <a href="/">
          <script id="img">
          document.write('<img src="/images/isucon-bank.png" alt="いすこん銀行 オンラインバンキングサービス">');
          script = document.getElementById('img');
          script.parentNode.removeChild(script);
          </script>
        </a>
      </h1>
      <div id="be-careful-phising" class="panel panel-danger">
  <div class="panel-heading">
    <span class="hikaru-mozi">偽画面にご注意ください！</span>
  </div>
  <div class="panel-body">
    <p>偽のログイン画面を表示しお客様の情報を盗み取ろうとする犯罪が多発しています。</p>
    <p>ログイン直後にダウンロード中や、見知らぬウィンドウが開いた場合、<br>すでにウィルスに感染している場合がございます。即座に取引を中止してください。</p>
    <p>また、残高照会のみなど、必要のない場面で乱数表の入力を求められても、<br>絶対に入力しないでください。</p>
  </div>
</div>

<div class="page-header">
  <h1>ログイン</h1>
</div>


<div class="container">
  <form class="form-horizontal" role="form" action="/login" method="POST">
    <div class="form-group">
      <label for="input-username" class="col-sm-3 control-label">お客様ご契約ID</label>
      <div class="col-sm-9">
        <input id="input-username" type="text" class="form-control" placeholder="半角英数字" name="login">
      </div>
    </div>
    <div class="form-group">
      <label for="input-password" class="col-sm-3 control-label">パスワード</label>
      <div class="col-sm-9">
        <input type="password" class="form-control" id="input-password" name="password" placeholder="半角英数字・記号（２文字以上）">
      </div>
    </div>
    <div class="form-group">
      <div class="col-sm-offset-3 col-sm-9">
        <button type="submit" class="btn btn-primary btn-lg btn-block">ログイン</button>
      </div>
    </div>
  </form>
</div>

    </div>

  </body>
</html>
		`))
	}
}

func postLogin(c *gin.Context) {
	user, err := attemptLogin(c.Request)

	if err != nil || user == nil {
		switch err {
		case ErrBannedIP:
			c.Redirect(http.StatusMovedPermanently, "/?err=banned")
		case ErrLockedUser:
			c.Redirect(http.StatusMovedPermanently, "/?err=locked")
		default:
			c.Redirect(http.StatusMovedPermanently, "/?err=wrong")
		}

		return
	}
	session, _ := store.Get(c.Request, "user_id")
	session.Values["user_id"] = user.ID
	session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusOK, "/mypage")
}

func getMypage(c *gin.Context) {
	var currentUser *User
	session, _ := store.Get(c.Request, "user_id")

	if userId, ok := session.Values["user_id"]; ok {
		currentUser = getCurrentUser(userId)
	} else {
		currentUser = nil
	}

	if currentUser == nil {
		c.Redirect(http.StatusMovedPermanently, "/?err=invalid")
		return
	}

	currentUser.getLastLogin()

	format := `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <script id="css">
    document.write('<link rel="stylesheet" href="/stylesheets/bootstrap.min.css"> <link rel="stylesheet" href="/stylesheets/bootflat.min.css"> <link rel="stylesheet" href="/stylesheets/isucon-bank.css">');
    script = document.getElementById('css');
    script.parentNode.removeChild(script);
    </script>
    <title>isucon4</title>
  </head>
  <body>
    <div class="container">
      <h1 id="topbar">
        <a href="/">
          <script id="img">
          document.write('<img src="/images/isucon-bank.png" alt="いすこん銀行 オンラインバンキングサービス">');
          script = document.getElementById('img');
          script.parentNode.removeChild(script);
          </script>
        </a>
      </h1>
      <div class="alert alert-success" role="alert">
        ログインに成功しました。<br>
        未読のお知らせが０件、残っています。
      </div>

      <dl class="dl-horizontal">
        <dt>前回ログイン</dt>
        <dd id="last-logined-at">%s</dd>
        <dt>最終ログインIPアドレス</dt>
        <dd id="last-logined-ip">%s</dd>
      </dl>

      <div class="panel panel-default">
        <div class="panel-heading">
          お客様ご契約ID：%s 様の代表口座
        </div>
        <div class="panel-body">
          <div class="row">
            <div class="col-sm-4">
              普通預金<br>
              <small>東京支店　1111111111</small><br>
            </div>
            <div class="col-sm-4">
              <p id="zandaka" class="text-right">
                ―――円
              </p>
            </div>

            <div class="col-sm-4">
              <p>
                <a class="btn btn-success btn-block">入出金明細を表示</a>
                <a class="btn btn-default btn-block">振込・振替はこちらから</a>
              </p>
            </div>

            <div class="col-sm-12">
              <a class="btn btn-link btn-block">定期預金・住宅ローンのお申込みはこちら</a>
            </div>
          </div>
        </div>
      </div>
    </div>

  </body>
</html>
	`
	c.Data(200, "text/html", []byte(
		fmt.Sprintf(
			format,
			currentUser.LastLogin.CreatedAt.Format("2006-01-02 15:04:05"),
			currentUser.LastLogin.IP,
			currentUser.LastLogin.Login,
		),
	))
	//c.HTML(http.StatusOK, "template/mypage.tmpl", gin.H{
	//	"LastLogin" : currentUser,
	//})
}

func getReport(c *gin.Context) {
	c.JSON(200, map[string][]string{
		"banned_ips":   bannedIPs(),
		"locked_users": lockedUsers(),
	})
}

func main() {
	runtime.GOMAXPROCS(4)

	r := gin.New()

	r.GET("/", getIndex)
	r.POST("/login", postLogin)
	r.GET("/mypage", getMypage)
	r.GET("/report", getReport)

	//initStaticFiles("../public")

/*
	m := martini.Classic()

	store := sessions.NewCookieStore([]byte("secret-isucon"))
	m.Use(sessions.Sessions("isucon_go_session", store))

	m.Use(martini.Static("../public"))
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Get("/", func(r render.Render, session sessions.Session) {
		r.HTML(200, "index", map[string]string{"Flash": getFlash(session, "notice")})
	})

	m.Post("/login", func(req *http.Request, r render.Render, session sessions.Session) {}
		user, err := attemptLogin(req)

		notice := ""
		if err != nil || user == nil {
			switch err {
			case ErrBannedIP:
				notice = "You're banned."
			case ErrLockedUser:
				notice = "This account is locked."
			default:
				notice = "Wrong username or password"
			}

			session.Set("notice", notice)
			r.Redirect("/")
			return
		}

		session.Set("user_id", strconv.Itoa(user.ID))
		r.Redirect("/mypage")
	})

	m.Get("/mypage", func(r render.Render, session sessions.Session) {
		currentUser := getCurrentUser(session.Get("user_id"))

		if currentUser == nil {
			session.Set("notice", "You must be logged in")
			r.Redirect("/")
			return
		}

		currentUser.getLastLogin()
		r.HTML(200, "mypage", currentUser)
	})

	m.Get("/report", func(r render.Render) {
		r.JSON(200, map[string][]string{
			"banned_ips":   bannedIPs(),
			"locked_users": lockedUsers(),
		})
	})
*/
	//http.ListenAndServe(":8080", m)
	r.Run(":8080")

}
