<html>
  <head>
    <title>Got Scripts?</title>
    <style type="text/css">
      html {
        font-family: arial,sans-serif;
	font-size: small;
      }
      ul {
        list-style-type: none;
      }
      li {
        display: inline;
	padding-right: 10px;
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
      }
      .alert {
        width: 100%;
        color: red;
	font-weight: bold;
      }
      .main {
        padding: 0 40px 0 40px;
      }
    </style>
  </head>
  <body>
    <ul>
      <li><a href="/">home</a></li>
      <li><a href="/scripts">scripts</a></li>
      <li><a href="/add">+add</a></li>
    </ul>
    {{#alert}}<center class="alert">{{alert}}</center>{{/alert}}
    {{{content}}}
    {{> _pager.m}}
  </body>
</html>
