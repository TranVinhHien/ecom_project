import { create } from 'zustand';

interface ChatStore {
  isOpen: boolean;
  pendingMessage: string | null;
  productKey: string | null;
  
  openChatWithMessage: (message: string, productKey?: string) => void;
  setIsOpen: (isOpen: boolean) => void;
  clearPendingMessage: () => void;
}

export const useChatStore = create<ChatStore>((set) => ({
  isOpen: false,
  pendingMessage: null, 
  productKey: null,
  
  openChatWithMessage: (message: string, productKey?: string) => {
    set({ 
      isOpen: true, 
      pendingMessage: message,
      productKey: productKey || null
    });
  },
  
  setIsOpen: (isOpen: boolean) => {
    set({ isOpen });
  },
  
  clearPendingMessage: () => {
    set({ pendingMessage: null, productKey: null });
  },
}));
