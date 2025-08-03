import { createRouter, RouterProvider } from "@tanstack/react-router";
import React from "react";
import ReactDOM from "react-dom/client";
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { routeTree } from './routes/__root';
import { IconContext } from "react-icons/lib";

const queryClient = new QueryClient();

const router = createRouter({ routeTree });

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
    <React.StrictMode>
        <IconContext.Provider  value={{ style: { width: '24px', height: '24px' } }}>
            <QueryClientProvider client={queryClient}>
                <RouterProvider router={router} />
            </QueryClientProvider>
        </IconContext.Provider>
    </React.StrictMode>,
);
