package net

var head = `
<!doctype html>
<html lang="zh-cmn-Hans">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="description" content="">
    <meta name="author" content="">

    <title>`

var middle = `</title>

    <!-- Bootstrap core CSS -->
    <link href="https://cdn.bootcdn.net/ajax/libs/twitter-bootstrap/4.5.3/css/bootstrap.min.css" rel="stylesheet">

  </head>
  <body>

    <main role="main" class="container">
      <div class="jumbotron">
        <h1>
                 <form action="`

var tail = `" method="post" enctype="multipart/form-data">
                         <input type="file" name="file"><br>
                         <input type="submit">
                 </form>
	</h1>
      </div>
    </main>

    <script src="https://code.jquery.com/jquery-3.2.1.slim.min.js"></script>
    <script src="https://cdn.bootcdn.net/ajax/libs/twitter-bootstrap/4.5.3/js/bootstrap.min.js"></script>
  </body>
</html>
`

func GetUploadPage(title, path string) string {
	return head + title + middle + path + tail
}
