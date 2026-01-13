import { defineStore } from 'pinia'
import { getChatHistory } from '../api/chat'

export const useChatStore = defineStore('chat', {
  state: () => ({
    currentChat: null as any, // 当前选中的好友
    messages: [] as any[]     // 当前聊天记录
  }),
  actions: {
    // 选中好友，并拉取历史记录
    async selectFriend(friend: any) {
      this.currentChat = friend
      this.messages = [] // 切换时先清空，防止闪烁
      try {
        const res: any = await getChatHistory({ target_id: friend.id })
        this.messages = res.reverse()
      } catch (e) {
        console.error(e)
      }
    },
    // 接收到新消息 (不管是别人发的，还是自己发的)
    addMessage(msg: any) {
      // 只有当消息属于当前聊天对象，或者我们在看群聊时，才推入数组
      if (!this.currentChat) return

      // 检查是否是重复消息（乐观更新和服务器推送的重复）
      const isDuplicate = this.messages.some(existingMsg =>
        existingMsg.content === msg.content &&
        existingMsg.from_user_id === msg.from_user_id &&
        existingMsg.to_user_id === msg.to_user_id 
      )

      if (!isDuplicate) {
        this.messages.push(msg)
      }
    }
  }
})