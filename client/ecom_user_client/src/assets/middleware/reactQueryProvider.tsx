"use client"
import { useToast } from "@/hooks/use-toast";
import { QueryCache, QueryClient, QueryClientProvider } from "@tanstack/react-query"
export default function TanstackQueryProvider({ children }: { children: React.ReactNode }) {
    const { toast } = useToast()

    const queryClient = new QueryClient({
        queryCache: new QueryCache({
            onError: (error: any) => {
                toast({
                    variant: "destructive",
                    title: error.message,
                    description: error?.response?.data?.message,
                })
            },
        }),
        defaultOptions: {
            mutations: {
                onError: (error: any) => {
                    toast({
                        variant: "destructive",
                        title: error.message,
                        description: error?.response?.data?.message,
                    })
                },
            },
        },
    });

    return (
        <QueryClientProvider client={queryClient}>
            {children}
        </QueryClientProvider>
    )
}
