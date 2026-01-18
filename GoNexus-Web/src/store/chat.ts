import { defineStore } from 'pinia'
import { getChatHistory, readMessage } from '../api/chat'

export const useChatStore = defineStore('chat', {
  state: () => ({
    currentChat: null as any, // 当前选中的好友
    messages: [] as any[],     // 当前聊天记录
    unreadMap: {} as Record<string, number> // 未读消息计数 key: 'friend_1' | 'group_2'
  }),
  actions: {
    async selectFriend(friend: any) {
      this.currentChat = friend
      this.messages = []
      // 清除未读
      const key = friend.isGroup ? `group_${friend.id}` : `friend_${friend.id}`
      this.unreadMap[key] = 0
      
      try {
        const chat_type = friend.isGroup ? 2 : 1
        
        // 标记已读
        readMessage({
            target_id: friend.id,
            chat_type
        }).catch(err => console.error("标记已读失败", err))

        const res: any = await getChatHistory({
          target_id: friend.id,
          chat_type
        })
        this.messages = res.reverse()
      } catch (e) {
        console.error(e)
      }
    },
    addMessage(msg: any) {
      // 获取当前用户ID (用于判断是否是自己发的消息)
      const myId = parseInt(localStorage.getItem('user_id') || '0')
      
      // 1. 解析消息归属
      let targetKey = ''
      let isForCurrentChat = false

      if (msg.chat_type === 2) {
         // 群聊
         const groupId = msg.to_user_id ?? msg.group_id ?? msg.target_id
         targetKey = `group_${groupId}`
         
         // 判断是否当前窗口
         if (this.currentChat && this.currentChat.isGroup && String(this.currentChat.id) === String(groupId)) {
             isForCurrentChat = true
         }
      } else {
         // 单聊
         // 如果是别人发给我的，归属是 msg.from_user_id
         // 如果是我发给别人的(多端同步)，归属是 msg.to_user_id
         const otherId = (msg.from_user_id === myId) ? msg.to_user_id : msg.from_user_id
         targetKey = `friend_${otherId}`

         // 判断是否当前窗口
         if (this.currentChat && !this.currentChat.isGroup && String(this.currentChat.id) === String(otherId)) {
             isForCurrentChat = true
         }
      }

      // 2. 如果是当前窗口，追加消息
      if (isForCurrentChat) {
        const isDuplicate = this.messages.some(existingMsg =>
            existingMsg.content === msg.content &&
            existingMsg.from_user_id === msg.from_user_id &&
            existingMsg.to_user_id === msg.to_user_id 
        )
        if (!isDuplicate) {
            this.messages.push(msg)
            
            // 如果是对方发来的消息，且当前窗口打开，立即标记已读
            if (msg.from_user_id !== myId) {
                const chat_type = msg.chat_type // 1 or 2
                // 单聊 target_id 是 from_user_id (对方), 群聊 target_id 是 group_id
                const target_id = (chat_type === 2) ? (msg.group_id ?? msg.target_id) : msg.from_user_id
                
                readMessage({
                    target_id: target_id,
                    chat_type: chat_type
                }).catch(e => console.error("自动标记已读失败", e))
            }
        }
      } else {
        // 3. 如果不是当前窗口，且不是自己发的，增加未读计数
        if (msg.from_user_id !== myId) {
            this.unreadMap[targetKey] = (this.unreadMap[targetKey] || 0) + 1
        }
      }
    },
    // 处理撤回消息
    handleRevokeMessage(msgId: number, content: string = '此消息已被撤回') {
      const index = this.messages.findIndex(m => m.id === msgId || m.msg_id === msgId)
      if (index !== -1) {
        // 修改消息类型为系统消息或特殊撤回类型
        // 这里我们简单将其改为 TypeSystem (6) 并更新内容
        this.messages[index].type = 6
        this.messages[index].msg_type = 6
        this.messages[index].content = content
      }
    },
    // 清除当前会话
    clearCurrentChat() {
      this.currentChat = null
      this.messages = []
    }
  }
})
