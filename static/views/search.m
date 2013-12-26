    <div class="main">
      <h1>Search Results</h1>
      <h3>"{{query}}"</h3>
      <h5>{{total}} record(s) found</h5>
      {{#info}}
        <h3>Problem while searching {{.}}</h3>
      {{/info}}
      <ul class="results">
      {{#results}}
        <li>{{> _script.m}}</li>
      {{/results}}
      </ul>
    </div>
