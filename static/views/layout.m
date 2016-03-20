<html>
  <head>
    <meta name="viewport" content="initial-scale=1">
    <meta name="description" content="The screenplay search engine to help you find the script you're looking for. 1000s of screenplays in pdf format.">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <title>Scridx - A Screenplay Index</title>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" />
    <style>
      .center {
        margin-top: 100px;
      }
      .t-ico {
	padding: 0;
      }
      .t-ico img {
        width: 16px;
        height: auto;
      }
    </style>
  </head>
  <body>
    <nav class="navbar navbar-default">
      <div class="container">
        <div class="navbar-header">
          <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
            <span class="sr-only">Toggle navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          <a class="navbar-brand" href="/">Scridx</a>
        </div>
        <div id="navbar" class="collapse navbar-collapse">
          <ul class="nav navbar-nav">
            <li><a href="/trending">Trending</a></li>
            <li><a href="/random">Random</a></li>
            <li>
              <a href="/add"><span class="glyphicon glyphicon-plus" aria-hidden="true"></span> Add</a>
            </li>
            <li>
              <a href="/"><span class="glyphicon glyphicon-search" aria-hidden="true"></span> Search</a>
            </li>
          </ul>
          <ul class="nav navbar-nav navbar-right">
            <li class="t-ico"><a href="https://twitter.com/scridx"><img src="/static/image/t.png"/></a></li>
          </ul>
        </div><!--/.nav-collapse -->

      </div>
    </nav>

    <div class="container">
      {{#alert}}
      <div class="row">
        <div class="col-md-12"> 
          <center class="alert">{{alert}}</center>
        </div>
      </div>
      {{/alert}}

      <div class="row">
        <div class="col-md-10 col-md-offset-1">
          {{{content}}}
        </div>
      </div>

      <div class="row">
        <div class="col-md-10 col-md-offset-1">
          {{> _pager.m}}
        </div>
      </div>
    </div>


    <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.2.2/jquery.min.js"></script>
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js"></script>
    <script>
      function shareLink(obj) {
	href = obj.getAttribute("href");
	input = obj.parentNode.parentNode.getElementsByTagName("input")[0];
	input.style.display = 'inline';
	input.value = location.protocol + '//' + location.hostname + href;
	input.select();
	return false;
      };
    </script>
    <script>
      (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
      (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
      m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
      })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

      ga('create', 'UA-41003159-1', 'scridx.com');
      ga('send', 'pageview');
    </script>
  </body>
</html>
