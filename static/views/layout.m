<html>
  <head>
    <meta name="viewport" content="initial-scale=1">
    <title>Got Scripts?</title>
    <style type="text/css">
      html {
        font-family: arial,sans-serif;
	font-size: small;
        color: #404040;
      }
      body {
	min-width: 320px;
        max-width: 960px;
	margin: 0 auto;
      }
      button {
	background: #f1f1f1;
	border: 1px solid #dcdcdc;
        padding:5px;
      }
      ul {
        list-style-type: none;
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
        padding: 10% 20% 10% 20%;
      }
      .center-form {
        text-align: center;
        padding: 5% 20% 5% 20%;
      }
      .input {
        width: 100%;
	padding: 5px;
	font-weight: bold;
	border: 1px solid #4285f4;
      }
      .nav {
        padding-top: 10px;
      }
      .nav a {
        text-decoration: none;
      }
      .nav a:hover {
        text-decoration: underline;
      }
      .alert {
        width: 100%;
        color: red;
	font-weight: bold;
      }
      .main {
        padding: 0 40px 0 40px;
      }
      .results {
        padding: 0;
      }
      .results li {
        display: block;
        padding-bottom: 10px;
      }
      .domain {
        color: #006621;
      }
      .link input {
	padding: 0;
        margin: 0;
        width: 200px;
	border: 1px solid #222;
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
    </style>
  </head>
  <body>
    <ul class="nav">
      <li><a href="/">home</a></li>
      <li><a href="/scripts">scripts</a></li>
      <li><a href="/add">+add</a></li>
      <li><a href="/"><img src="/static/image/s.png"/></a></li>
    </ul>
    {{#alert}}<center class="alert">{{alert}}</center>{{/alert}}
    {{{content}}}
    {{> _pager.m}}
    <script>
      function shareLink(obj) {
	href = obj.getAttribute("href");
	input = obj.parentNode.getElementsByTagName("input")[0];
	input.style.display = 'inline-block';
	input.value = location.protocol + '//' + location.hostname + href;
	input.select();
	return false;
      };
    </script>
  </body>
</html>
