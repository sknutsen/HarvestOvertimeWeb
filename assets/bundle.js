function hoursAfterRequest() {
  enableButton("form__submit");
}

function hoursBeforeRequest() {
  disableButton("form__submit");
}

function disableButton(buttonId) {
  let button = document.getElementById(buttonId);

  if (button.getAttribute("disabled") !== "true") {
    button.setAttribute("disabled", "true");
  }
}

function enableButton(buttonId) {
  let button = document.getElementById(buttonId);

  if (button.getAttribute("disabled") === "true") {
    button.removeAttribute("disabled");
  }
}
