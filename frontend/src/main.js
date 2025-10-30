/* eslint-disable space-before-function-paren */
/* eslint-disable no-multi-spaces */
/* eslint-disable semi */
/* eslint-disable no-undef */
/* eslint-disable indent */
/* eslint-disable quotes */
import './assets/main.css'
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import store from './store'
import './assets/styles/github-markdown.css'
import './assets/styles/github-highlight.css'
import headImage from './assets/my_headportrait.jpg'
import {
  faHouse,
  faBlog,
  faMagnifyingGlass,
  faQuestion,
  faDiagramProject,
  faFlask,
  faPenToSquare,
  faComment,
  faCommentSlash,
  faRightToBracket,
  faRegistered,
  faArrowRotateBack,
  faEllipsis,
  faBoxOpen,
  faDice,
  faHeading,
  faList,
  faPaperPlane,
  faFilm,
  faBars,
  faBookmark,
  faCheckCircle,
  faBold,
  faItalic,
  faCode,
  faFileCode,
  faQuoteLeft,
  faTrash,
  faCalendar,
  faEye,
  faEyeSlash,
  faCommentDots,
  faImage,
  faEdit,
  faHome,
  faArrowLeft,
  faSyncAlt,
  faReply,
  faClock,
  faImages,
  faSun,
  faMoon,
  faShare,
  faArchive,
  faRss,
  faArrowUp,
  faFire,
  faTag,
  faLightbulb,
  faTimes,
  faBolt,
  faBookOpen,
  faChevronLeft,
  faChevronRight,
  faStrikethrough,
  faListUl,
  faListOl,
  faMinus,
  faTable,
  faLink,
  faHeart,
  faExclamationTriangle,
  faRedo,
  faUpload,
  faGear,
  faLocationDot,
  faNewspaper,
  faExternalLinkAlt,
  faGlobe,
  faBook,
  faDownload,
  faEnvelope,
  faPhone,
  faMapMarkerAlt,
  faStar,
  faCog,
  faUser,
  faInfoCircle,
  faWarning,
  faTimesCircle,
  faCheckCircle as faCheckCircleAlt,
  faQuestionCircle,
  faAddressBook,
  faVideo,
  faVolumeUp,
  faFile,
  faFolder,
  faTerminal,
  faDatabase,
  faServer,
  faCloud,
  faChalkboard,
  faCloudUploadAlt,
  faFilePowerpoint
} from '@fortawesome/free-solid-svg-icons'

import {
  faGithub,
  faGitlab,
  faBitbucket,
  faWeixin,
  faWeibo,
  faQq,
  faBilibili,
  faYoutube,
  faTwitter,
  faFacebook,
  faInstagram,
  faLinkedin,
  faGit
} from '@fortawesome/free-brands-svg-icons'

library.add(
  faHouse,
  faBlog,
  faMagnifyingGlass,
  faQuestion,
  faDiagramProject,
  faFlask,
  faPenToSquare,
  faComment,
  faCommentSlash,
  faRightToBracket,
  faRegistered,
  faArrowRotateBack,
  faEllipsis,
  faBoxOpen,
  faDice,
  faHeading,
  faList,
  faPaperPlane,
  faFilm,
  faBars,
  faBookmark,
  faCheckCircle,
  faBold,
  faItalic,
  faCode,
  faFileCode,
  faQuoteLeft,
  faTrash,
  faCalendar,
  faEye,
  faEyeSlash,
  faCommentDots,
  faImage,
  faEdit,
  faHome,
  faArrowLeft,
  faSyncAlt,
  faReply,
  faClock,
  faImages,
  faSun,
  faMoon,
  faShare,
  faArchive,
  faRss,
  faArrowUp,
  faFire,
  faTag,
  faLightbulb,
  faTimes,
  faBolt,
  faBookOpen,
  faChevronLeft,
  faChevronRight,
  faStrikethrough,
  faListUl,
  faListOl,
  faMinus,
  faTable,
  faLink,
  faHeart,
  faExclamationTriangle,
  faRedo,
  faUpload,
  faGear,
  faLocationDot,
  faNewspaper,
  faExternalLinkAlt,
  faGlobe,
  faBook,
  faDownload,
  faEnvelope,
  faPhone,
  faMapMarkerAlt,
  faStar,
  faCog,
  faUser,
  faInfoCircle,
  faWarning,
  faTimesCircle,
  faCheckCircleAlt,
  faQuestionCircle,
  faAddressBook,
  faVideo,
  faVolumeUp,
  faFile,
  faFolder,
  faTerminal,
  faDatabase,
  faServer,
  faCloud,
  faChalkboard,
  faCloudUploadAlt,
  faFilePowerpoint,
  faGithub,
  faGitlab,
  faBitbucket,
  faWeixin,
  faWeibo,
  faQq,
  faBilibili,
  faYoutube,
  faTwitter,
  faFacebook,
  faInstagram,
  faLinkedin,
  faGit
)

const app = createApp(App)
app.use(router)
app.use(store)
app.component('FontAwesomeIcon', FontAwesomeIcon)
app.mount('#app')

// 设置浏览器标签栏标题与图标
try {
  document.title = '山角函兽的小窝'
  const setFavicon = (href) => {
    let link = document.querySelector("link[rel='icon']")
    if (!link) {
      link = document.createElement('link')
      link.rel = 'icon'
      document.head.appendChild(link)
    }
    link.href = href
  }
  // 生成带圆角的 favicon
  const createRoundedFavicon = (src, size = 64, radius = 12) => new Promise((resolve, reject) => {
    const img = new Image()
    img.crossOrigin = 'anonymous'
    img.onload = () => {
      try {
        const canvas = document.createElement('canvas')
        canvas.width = size
        canvas.height = size
        const ctx = canvas.getContext('2d')

        const r = Math.min(radius, size / 2)
        const w = size
        const h = size
        ctx.beginPath()
        ctx.moveTo(r, 0)
        ctx.arcTo(w, 0, w, h, r)
        ctx.arcTo(w, h, 0, h, r)
        ctx.arcTo(0, h, 0, 0, r)
        ctx.arcTo(0, 0, w, 0, r)
        ctx.closePath()
        ctx.clip()

        // 等比裁剪
        const minSide = Math.min(img.width, img.height)
        const sx = (img.width - minSide) / 2
        const sy = (img.height - minSide) / 2
        ctx.drawImage(img, sx, sy, minSide, minSide, 0, 0, size, size)
        resolve(canvas.toDataURL('image/png'))
      } catch (e) {
        reject(e)
      }
    }
    img.onerror = () => reject(new Error('favicon image load error'))
    img.src = src
  })

  createRoundedFavicon(headImage, 64, 12)
    .then(url => setFavicon(url))
    .catch(() => setFavicon(headImage))
} catch (e) {
  // 忽略在非浏览器环境下的错误
}
