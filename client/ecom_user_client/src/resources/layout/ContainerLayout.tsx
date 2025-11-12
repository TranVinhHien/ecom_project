"use client"
import { useRouter } from "@/i18n/routing"
import { CSSProperties, ReactNode, useEffect, useState } from "react";
import Header from "./header";
import Sidebar from "./sidebar";
import { Toaster } from "@/components/ui/toaster";
import Footer from "./footer";
import ChatBox from "@/components/ChatBox";
import { useTokenRefresh } from "@/hooks/useTokenRefresh";

const Layout = ({ children }:
    {
        children: ReactNode
    }) => {

    // Initialize auto token refresh
    useTokenRefresh();

    return (
        <div>
            <div>
                    <Toaster />
                        <Header onClick={()=>{}} />
                    <div className="mt-4"></div>
                    <div className="relative px-14" style={{ flex: 1 }}>
                            {children}
                    </div>
                    
                    <Footer />
                    <ChatBox />
            </div>
        </div>
    );
};

export default Layout;
