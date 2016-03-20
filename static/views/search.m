<div class="page-header">
  <form action="/search" method="get">
    <div class="input-group col-md-6">
      <input class="form-control form-inline" type="text" name="q" id="q" value="{{query}}"/>
      <span class="input-group-btn">
        <button class="btn btn-default form-inline">Search</button>
      </span>
    </div>
  </form>
  <h5>{{total}} results</h5>
</div>

  {{#info}}
    <h3>Problem while searching {{.}}</h3>
  {{/info}}

  {{#results}}
  <div class="row">
    <div class="col-md-12">
      {{> _script.m}}
    </div>
  </div>
  {{/results}}
