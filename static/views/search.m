    <div class="main">
      <h1>Search Results</h1>
      <h3>"{{query}}"</h3>
      <h5>{{total}} record(s) found</h5>
      {{#info}}
        <h3>Problem while searching {{.}}</h3>
      {{/info}}
      {{#results}}
        <div><a href="{{Url}}">{{Title}}</a> <a href="/s/{{Short}}" style="font-size:0.8em">short</a></div>
      {{/results}}
    </div>
