const API_BASE = (() => {
  if (window.location.protocol === "file:") {
    return "http://localhost:8080";
  }
  if (window.location.port === "3000") {
    return `${window.location.origin}/api`;
  }
  return "http://localhost:8080";
})();
