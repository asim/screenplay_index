<center class="center">
  <h1>Add Script</h1>

  <form action="/add" method="post">
    <div class="input-group col-md-8">
      <br>
      <label>Title</label>
      <input class="form-control" type="text" name="title"/>
    </div>

    <div class="input-group col-md-8">
      <br>
      <label>URL (pdf only)</label>
      <input class="form-control" type="text" name="url"/>
    <div>

    <div class="input-group col-md-8">
      <br>
      <label>Captcha</label>
      <p><img src="/captcha/{{captcha}}.png"/></p>
    </div>

    <div class="input-group col-md-8">
      <input class="form-control" type="text" name="captcha"/>
      <input type="hidden" name="_captchaId" value="{{captcha}}"/>
    </div>

    <div class="form-group">
      <br>
      <button class="btn btn-default">Add</button>
    </div>
  </form>
</center>
