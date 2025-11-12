"use client";

import { useState, useEffect, useRef } from "react";
import { MessageCircle, X, Send, Loader2, ThumbsUp, ThumbsDown, Router } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";
import { getCookieValues } from "@/assets/helpers/cookies";
import { ACCESS_TOKEN, MESSAGES_STORAGE_KEY, SESSION_STORAGE_KEY } from "@/assets/configs/request";
import { useToast } from "@/hooks/use-toast";
import ProductCarousel from "@/components/ProductCarousel";
import API from "@/assets/configs/api";
import { generateComplaintUrl } from "@/lib/complaintUtils";
import { useRouter } from "@/i18n/routing";

interface Message {
  id: string;
  text: string;
  isUser: boolean;
  timestamp: Date;
  products?: MatchedProduct[];
  event_id?: string;
  userPrompt?: string;
  rating?: number | null;
}

interface MatchedProduct {
  brand: string;
  category: string;
  product: {
    id: string;
    image: string;
    key: string;
    name: string;
    product_is_permission_check: boolean;
    product_is_permission_return: boolean;
    short_description: string;
  };
  sku: Array<{
    id: string;
    price: number;
    quantity: number;
    sku_name: string;
  }>;
  similarity_score: number;
}



// Export function to clear chat history (used on logout)
export const clearChatHistory = () => {
  localStorage.removeItem(SESSION_STORAGE_KEY);
  localStorage.removeItem(MESSAGES_STORAGE_KEY);
};

export default function ChatBox() {
  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputMessage, setInputMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [sessionId, setSessionId] = useState<string | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const { toast } = useToast();
  const isInitialized = useRef(false);
  const router = useRouter();
  // Load chat history from localStorage on mount
  useEffect(() => {
    if (!isInitialized.current) {
      isInitialized.current = true;
      const savedSessionId = localStorage.getItem(SESSION_STORAGE_KEY);
      const savedMessages = localStorage.getItem(MESSAGES_STORAGE_KEY);

      if (savedSessionId) {
        setSessionId(savedSessionId);
      }

      if (savedMessages) {
        try {
          const parsedMessages = JSON.parse(savedMessages);
          // Convert timestamp strings back to Date objects
          const messagesWithDates = parsedMessages.map((msg: any) => ({
            ...msg,
            timestamp: new Date(msg.timestamp),
          }));
          setMessages(messagesWithDates);
        } catch (error) {
          console.error("Failed to parse saved messages:", error);
        }
      }
    }
  }, []);

  // Save messages to localStorage whenever they change
  useEffect(() => {
    if (messages.length > 0) {
      localStorage.setItem(MESSAGES_STORAGE_KEY, JSON.stringify(messages));
    }
  }, [messages]);

  // Save sessionId to localStorage whenever it changes
  useEffect(() => {
    if (sessionId) {
      localStorage.setItem(SESSION_STORAGE_KEY, sessionId);
    }
  }, [sessionId]);

  // Auto scroll to bottom
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // Create session when opening chat
  useEffect(() => {
    if (isOpen && !sessionId) {
      createSession();
    }
  }, [isOpen, sessionId]);

  // Remove cleanup on unmount - keep chat history
  // Only cleared on logout

  const createSession = async () => {
    try {
      const token = getCookieValues<string>(ACCESS_TOKEN);
      if (!token) {
        toast({
          title: "Lỗi xác thực",
          description: "Vui lòng đăng nhập để sử dụng chat",
          variant: "destructive",
        });
        return;
      }

      const response = await fetch(`${API.base_agent}${API.agent.session}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          state: {
            lang: "VN",
          },
        }),
      });

      const data = await response.json();

      if (data.success && data.session_id) {
        setSessionId(data.session_id);
        
        // Only add welcome message if there are no existing messages
        if (messages.length === 0) {
          setMessages([
            {
              id: Date.now().toString(),
              text: "Xin chào! Tôi là trợ lý mua sắm của bạn. Hãy cho tôi biết bạn muốn tìm sản phẩm gì nhé!",
              isUser: false,
              timestamp: new Date(),
            },
          ]);
        }
      } else {
        throw new Error("Không thể tạo phiên chat");
      }
    } catch (error) {
      console.error("Create session error:", error);
      toast({
        title: "Lỗi",
        description: "Không thể khởi tạo chat. Vui lòng thử lại.",
        variant: "destructive",
      });
    }
  };

  const sendMessage = async () => {
    if (!inputMessage.trim() || !sessionId) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      text: inputMessage,
      isUser: true,
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setInputMessage("");
    setIsLoading(true);

    try {
      const token = getCookieValues<string>(ACCESS_TOKEN);
      if (!token) {
        throw new Error("Token not found");
      }

      const response = await fetch(`${API.base_agent}${API.agent.message}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          message: inputMessage,
          session_id: sessionId,
        }),
      });

      const data = await response.json();

      if (data.success && data.response) {
        const botMessage: Message = {
          id: (Date.now() + 1).toString(),
          text: data.response.text,
          isUser: false,
          timestamp: new Date(),
          products: data.response.matched_products || undefined,
          event_id: data.event_id || (Date.now() + 1).toString(),
          userPrompt: inputMessage,
          rating: null,
        };
        
        setMessages((prev) => [...prev, botMessage]);
        if (data.response.category){
        const url = generateComplaintUrl(
          {
            category: data.response?.category,
            content: data.response?.context,
            
          }
        )
        router.push(url);
      } 
      } else {
        throw new Error(data.error || "Có lỗi xảy ra");
      }
    } catch (error: any) {
      console.error("Send message error:", error);
      toast({
        title: "Lỗi gửi tin nhắn",
        description: error.message || "Không thể gửi tin nhắn. Vui lòng thử lại.",
        variant: "destructive",
      });

      // Add error message
      setMessages((prev) => [
        ...prev,
        {
          id: (Date.now() + 1).toString(),
          text: "Xin lỗi, tôi không thể xử lý yêu cầu của bạn lúc này. Vui lòng thử lại sau.",
          isUser: false,
          timestamp: new Date(),
        },
      ]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  const toggleChat = () => {
    setIsOpen(!isOpen);
  };

  const handleRating = async (messageId: string, rating: number) => {
    if (!sessionId) return;

    // Find the message being rated
    const messageIndex = messages.findIndex((m) => m.id === messageId);
    if (messageIndex === -1) return;

    const message = messages[messageIndex];
    
    // Find the previous user message (the prompt for this response)
    let userPrompt = "";
    for (let i = messageIndex - 1; i >= 0; i--) {
      if (messages[i].isUser) {
        userPrompt = messages[i].text;
        break;
      }
    }

    try {
      const response = await fetch(`${API.base_analytics}${API.analytics.chatboxReview}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          session_id: sessionId,
          event_id: message.event_id || message.id,
          rating: rating,
          user_prompt: userPrompt,
          agent_response: message.text,
        }),
      });

      const data = await response.json();

      if (data.code === 200 && data.status === "success") {
        // Update the message with the rating
        setMessages((prev) =>
          prev.map((m) =>
            m.id === messageId ? { ...m, rating: rating } : m
          )
        );

        toast({
          title: "Thành công",
          description: data.result?.message || "Cảm ơn bạn đã đánh giá!",
        });
      } else {
        throw new Error(data.message || "Đánh giá thất bại");
      }
    } catch (error: any) {
      console.error("Rating error:", error);
      toast({
        title: "Lỗi",
        description: error.message || "Không thể gửi đánh giá. Vui lòng thử lại.",
        variant: "destructive",
      });
    }
  };

  return (
    <>
      {/* Chat Button */}
      {!isOpen && (
        <Button
          onClick={toggleChat}
          className="fixed bottom-6 right-6 h-14 w-14 rounded-full shadow-lg bg-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/.9)] z-50"
          size="icon"
        >
          <MessageCircle className="h-6 w-6" />
        </Button>
      )}

      {/* Chat Window */}
      {isOpen && (
        <Card className="fixed bottom-6 right-6 w-[400px] h-[600px] flex flex-col shadow-2xl z-50 overflow-hidden">
          {/* Header */}
          <div className="bg-[hsl(var(--primary))] text-white p-4 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-white/20 flex items-center justify-center">
                <MessageCircle className="h-5 w-5" />
              </div>
              {/* <div>
                <h3 className="font-semibold">Trợ lý mua sắm</h3>
                <p className="text-xs opacity-90">Luôn sẵn sàng hỗ trợ bạn</p>
              </div> */}
            </div>
            <Button
              variant="ghost"
              size="icon"
              onClick={toggleChat}
              className="text-white hover:bg-white/20"
            >
              <X className="h-5 w-5" />
            </Button>
          </div>

          {/* Messages */}
          <div className="flex-1 overflow-y-auto p-4 bg-gray-50 space-y-4">
            {messages.map((message) => (
              <div
                key={message.id}
                className={`flex ${message.isUser ? "justify-end" : "justify-start"}`}
              >
                <div
                  className={`max-w-[80%] rounded-lg p-3 ${
                    message.isUser
                      ? "bg-[hsl(var(--primary))] text-white"
                      : "bg-white border shadow-sm"
                  }`}
                >
                  <p className="text-sm whitespace-pre-wrap">{message.text}</p>
                  {message.products && message.products.length > 0 && (
                    <div className="mt-3">
                      <ProductCarousel products={message.products} />
                    </div>
                  )}
                  <p
                    className={`text-xs mt-2 ${
                      message.isUser ? "text-white/70" : "text-gray-500"
                    }`}
                  >
                    {message.timestamp.toLocaleTimeString("vi-VN", {
                      hour: "2-digit",
                      minute: "2-digit",
                    })}
                  </p>

                  {/* Like/Dislike buttons for agent messages */}
                  {!message.isUser && (
                    <div className="flex items-center gap-2 mt-2 pt-2 border-t border-gray-200">
                      <Button
                        variant="ghost"
                        size="sm"
                        className={`h-7 px-2 hover:bg-green-50 ${
                          message.rating === 1 ? "bg-green-100 text-green-600" : "text-gray-500"
                        }`}
                        onClick={() => handleRating(message.id, 1)}
                        disabled={message.rating !== null && message.rating !== 1}
                      >
                        <ThumbsUp className="h-3 w-3 mr-1" />
                        <span className="text-xs">Hữu ích</span>
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        className={`h-7 px-2 hover:bg-red-50 ${
                          message.rating === 0 ? "bg-red-100 text-red-600" : "text-gray-500"
                        }`}
                        onClick={() => handleRating(message.id, -1)}
                        disabled={message.rating !== null && message.rating !== -1}
                      >
                        <ThumbsDown className="h-3 w-3 mr-1" />
                        <span className="text-xs">Chưa tốt</span>
                      </Button>
                    </div>
                  )}
                </div>
              </div>
            ))}

            {isLoading && (
              <div className="flex justify-start">
                <div className="bg-white border shadow-sm rounded-lg p-3 flex items-center gap-2">
                  <Loader2 className="h-4 w-4 animate-spin text-[hsl(var(--primary))]" />
                  <div className="flex gap-1">
                    <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: "0ms" }}></span>
                    <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: "150ms" }}></span>
                    <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: "300ms" }}></span>
                  </div>
                </div>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>

          {/* Input */}
          <div className="p-4 bg-white border-t">
            <div className="flex gap-2">
              <Input
                value={inputMessage}
                onChange={(e) => setInputMessage(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="Nhập tin nhắn..."
                disabled={isLoading || !sessionId}
                className="flex-1"
              />
              <Button
                onClick={sendMessage}
                disabled={isLoading || !inputMessage.trim() || !sessionId}
                className="bg-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/.9)]"
                size="icon"
              >
                {isLoading ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <Send className="h-4 w-4" />
                )}
              </Button>
            </div>
          </div>
        </Card>
      )}
    </>
  );
}
