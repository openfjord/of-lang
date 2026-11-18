// The site is dark unless the reader chose the light theme.
(function () {
  var root = document.documentElement;
  try {
    if (localStorage.getItem("theme") === "light") root.dataset.theme = "light";
  } catch (error) {
    // No storage: the page stays dark.
  }
  document.addEventListener("DOMContentLoaded", function () {
    // On a small screen, the Menu button opens and closes the navigation.
    var menu = document.getElementById("menu");
    menu.addEventListener("click", function () {
      var open = root.classList.toggle("nav-open");
      menu.setAttribute("aria-expanded", open ? "true" : "false");
    });
    document.getElementById("theme").addEventListener("click", function () {
      var light = root.dataset.theme !== "light";
      if (light) root.dataset.theme = "light";
      else delete root.dataset.theme;
      try {
        localStorage.setItem("theme", light ? "light" : "dark");
      } catch (error) {
        // No storage: the choice lasts for this page only.
      }
    });
  });
})();
