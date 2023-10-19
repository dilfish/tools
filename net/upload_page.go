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
    <link href="https://cdn.dilfish.icu/302/bootstrap.css" rel="stylesheet">

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

    <main role="main" class="container">
        <div class="jumbotron">
            <h1>累积上传最多1G，单次最大10M</h1>
            <h1>curl -X POST -H "Content-Type: multipart/form-data" -F "file=@filename.fileext" https://dev.ug/upload</h1>
        </div>
    </main>

    <script src="https://cdn.dilfish.icu/302/jquery.js"></script>
    <script src="https://cdn.dilfish.icu/302/bootstrap.js"></script>
  </body>
</html>
`

func GetUploadPage(title, path string) string {
	return head + title + middle + path + tail
}
