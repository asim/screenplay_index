    <div class="center">
      <h1>Add Script</h1>
      <form action="/add" method="post">
	<p>Title</p>
        <input class="input" type="text" name="title"/>
	<p>Url (pdf only)</p>
        <input class="input" type="text" name="url"/>
 	<p>Captcha</p>
        <p><img src="/captcha/{{captcha}}.png"/></p>
       <input type="text" name="captcha"/>
       <input type="hidden" name="_captchaId" value="{{captcha}}"/>
        <p><button value="submit">Add</button>
      </form>
    </div>
