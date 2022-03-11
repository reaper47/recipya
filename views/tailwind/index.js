import picoModal from 'picomodal'
import Toastify from 'toastify-js'

const icons = {
  close: 'M6 18L18 6M6 6l12 12',
  hamburger: 'M4 6h16M4 12h16M4 18h16'
}

document
  .getElementById('nav__avatar-button')
  .addEventListener('mousedown', () => {
    document.getElementById('nav__avatar-content').classList.toggle('hidden')
  })

document.onclick = ({ target }) => {
  const id = target.id

  const avatarButton = document.getElementById('nav__avatar-button')
  const avatarMenuContent = document.getElementById('nav__avatar-content')

  if (![avatarMenuContent.id, avatarButton.id].includes(id)) {
    avatarMenuContent.classList.add('hidden')
  }
}

// Nav sidebar
document.getElementById('nav__menu-icon').addEventListener('mousedown', () => {
  const menuIconPath = document.getElementById('nav__menu-icon-path')
  if (menuIconPath.getAttribute('d') === icons.close) {
    menuIconPath.setAttribute('d', icons.hamburger)
  } else {
    menuIconPath.setAttribute('d', icons.close)
  }

  document.getElementById('nav__sidebar').classList.toggle('hidden')
})

const sidebar = document.getElementById('sidebar__recipes')
switch (document.location.pathname) {
  case '/':
  case '/recipes':
    sidebar.classList.add('border-l-4', 'border-red-600')
    break
  default:
    break
}

// Global functions
window.openModal = (event) => {
  event = event.target
  while (!event.id) {
    event = event.parentElement
  }

  window.modal = picoModal({
    content: document.getElementById(event.id.split('-open-button')[0])
  })
    .beforeShow((modal) =>
      modal.modalElem().children[0].classList.remove('hidden')
    )
    .beforeClose((modal) =>
      modal.modalElem().children[0].classList.add('hidden')
    )
    .show()
}

window.showSuccessToast = (text) =>
  Toastify({
    text,
    close: true,
    position: 'center',
    style: {
      background: '#27ae60'
    }
  }).showToast()

window.showErrorToast = (text) =>
  Toastify({
    text,
    close: true,
    position: 'center',
    style: {
      background: '#c0392b'
    }
  }).showToast()

// Storage
const toastMessage = window.sessionStorage.getItem('showSuccessToast')
if (toastMessage !== null) {
  window.showSuccessToast(toastMessage)
  window.sessionStorage.removeItem('showSuccessToast')
}
