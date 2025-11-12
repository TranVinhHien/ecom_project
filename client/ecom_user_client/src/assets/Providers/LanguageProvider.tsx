"use client";

import React from "react";
export default function LanguageProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  return <div >{children}</div>;
}
