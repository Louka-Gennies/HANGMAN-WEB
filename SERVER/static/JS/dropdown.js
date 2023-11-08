/* quand l'utilisateur clique sur le,
changer entre montrer et cacher le menu du dropdown */
function myFunction()
{
  document.getElementById("myDropdown").classList.toggle("show");
}

// fermer le menu quand l'utilisateur clique sur close
window.onclick = function (event)
{
  if (!event.target.matches('.dropbtn'))
  {
    var dropdowns = document.getElementsByClassName("dropdown-content");
    var i;
    for (i = 0; i < dropdowns.length; i++)
    {
      var openDropdown = dropdowns[i];
      if (openDropdown.classList.contains('show'))
      {
        openDropdown.classList.remove('show');
      }
    }
  }
}