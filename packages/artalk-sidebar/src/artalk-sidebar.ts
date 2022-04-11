import Artalk from 'artalk'
import ArtalkConfig from 'artalk/types/artalk-config'
import { UserData } from 'artalk/types/artalk-data'
import 'artalk/dist/Artalk.css'
import Sidebar from './sidebar'

class ArtalkSidebar extends Artalk {
  constructor(customConf: ArtalkConfig, user: UserData) {
    super(customConf)

    this.$root.style.display = 'none'
    this.ctx.user.data = user
    this.ctx.user.save()

    const sidebar = new Sidebar(this.ctx)
    document.body.appendChild(sidebar.$el)

    sidebar.show()
    console.log('hello artalk-sidebar')
  }
}

export default ArtalkSidebar
