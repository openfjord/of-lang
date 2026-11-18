// The site is dark unless the reader chose the light theme.
(function () {
  var root = document.documentElement;
  try {
    if (localStorage.getItem("theme") === "light") root.dataset.theme = "light";
  } catch (error) {
    // No storage: the page stays dark.
  }
  document.addEventListener("DOMContentLoaded", function () {
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
