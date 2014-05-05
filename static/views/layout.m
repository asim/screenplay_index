<html>
  <head>
    <meta name="viewport" content="initial-scale=1">
    <meta name="description" content="The screenplay search engine to help you find the script you're looking for. 1000s of screenplays in pdf format.">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <title>Scridx - A Screenplay Index</title>
    <style type="text/css">
      html {
        font-family: arial,sans-serif;
	font-size: small;
        color: #595959;
      }
      body {
	padding: 10px 10% 10px 10%;
	min-width: 280px;
        max-width: 960px;
	margin: 0 auto;
      }
      body a {
	color: #595959;
	font-size: 1.2em;
	font-weight: bold;
        text-decoration: none;
      }
      body a:hover {
	text-decoration: underline;
      }
      button {
	background: #f1f1f1;
	border: 1px solid #dcdcdc;
        padding:5px;
      }
      button:active {
        background: #E4E4E4;
      }
      button:hover {
        border: 1px solid #C3C3C3;
      }
      ul {
        list-style-type: none;
	padding: 0;
      }
      li {
        display: inline;
	padding-right: 10px;
      }
      ul img {
	vertical-align: bottom;
      }
      .center {
        text-align: center;
	padding-top: 10%;
      }
      .center .search {
        margin: 0 auto;
      }
      .input {
        width: 100%;
        max-width: 600px;
	padding: 5px;
	color: #595959;
	font-weight: bold;
	border: 1px solid #4285f4;
      }
      .nav {
        padding: 0;
	display: inline-block;
      }
      .nav a {
	font-size: small;
        text-decoration: none;
      }
      .nav a:hover {
        text-decoration: underline;
      }
      .alert {
        width: 100%;
        color: #dd4b39;
	font-weight: bold;
      }
      .meta {
        color: #535353;
      }
      .random-results {
	margin: 0 auto;
      }
      .random-results li h1 {
	margin-bottom: 0;
      }
      .results {
        padding: 0;
	max-width: 600px;
      }
      .results li {
        display: block;
        padding-bottom: 10px;
      }
      .right {
        float: right;
      }
      .doodle img {
        max-width: 100%;
        height: auto;
        width: auto\9;"
      }
      .domain {
        color: #858585;
	margin: 1px 0 1px 0;
      }
      .search {
        max-width: 600px;
      }
      .search div {
        overflow: hidden;
        max-width: 550px;
      }
      .search button {
        float: right;
      }
      .search *:focus {
        outline: none;
      }
      .t-ico {
	padding: 0;
      }
      .t-ico img {
        width: 16px;
        height: auto;
      }
      .share {
        display: inline;
        padding-left: 5px;
      }
      .share input {
	padding: 0;
        margin: 0;
        width: 200px;
	border: 1px solid #4285f4;
        background: #f1f1f1;
      }
      .share li {
        display: inline;
	padding: 0;
      }
    </style>
  </head>
  <body>
    <ul class="nav">
      <li class="home-ico"><a href="/"><img src="/static/image/h.png"/></a></li>
      <li><a href="/scripts">latest</a></li>
      <li><a href="/trending">trending</a></li>
      <li><a href="/random">random</a></li>
      <li><a href="/add">+add</a></li>
      <li class="search-ico"><a href="/"><img src="/static/image/s.png"/></a></li>
    </ul>
    <ul class="nav right">
      <li class="t-ico right"><a href="https://twitter.com/scridx"><img src="/static/image/t.png"/></a></li>
    </ul>
    {{#alert}}<center class="alert">{{alert}}</center>{{/alert}}
    {{{content}}}
    {{> _pager.m}}
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
