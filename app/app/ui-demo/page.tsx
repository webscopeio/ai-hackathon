"use client";

import { useState, useEffect, useRef } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";

export default function WebSocketDemo() {
  const [messages, setMessages] = useState<string[]>([]);
  const [inputValue, setInputValue] = useState("");
  const [isConnected, setIsConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const scrollAreaRef = useRef<HTMLDivElement>(null);

  const addMessage = (message: string) => {
    setMessages((prev) => [...prev, message]);
  };

  const handleOpen = () => {
    if (wsRef.current) return;

    try {
      console.log("Attempting to connect to WebSocket...");
      wsRef.current = new WebSocket("ws://localhost:8080/ws");

      wsRef.current.onopen = () => {
        console.log("WebSocket connection opened successfully");
        setIsConnected(true);
        addMessage("OPEN");
      };

      wsRef.current.onclose = (event) => {
        console.log("WebSocket connection closed:", event.code, event.reason);
        setIsConnected(false);
        addMessage(
          `CLOSE (Code: ${event.code}${
            event.reason ? `, Reason: ${event.reason}` : ""
          })`
        );
        wsRef.current = null;
      };

      wsRef.current.onmessage = (evt) => {
        console.log("Received message:", evt.data);
        addMessage(`RESPONSE: ${evt.data}`);
      };

      wsRef.current.onerror = (error) => {
        console.error("WebSocket error:", error);
        addMessage(
          `ERROR: Connection failed - check browser console for details`
        );
      };
    } catch (error: any) {
      console.error("WebSocket connection error:", error);
      addMessage(
        `ERROR: Failed to connect to WebSocket server - ${error.message}`
      );
    }
  };

  const handleClose = () => {
    if (!wsRef.current) return;
    wsRef.current.close();
  };

  const handleSend = () => {
    if (!wsRef.current) {
      try {
        console.log("Attempting to connect to WebSocket...");
        wsRef.current = new WebSocket("ws://localhost:8080/ws");

        wsRef.current.onopen = () => {
          console.log("WebSocket connection opened successfully");
          setIsConnected(true);
          addMessage("OPEN");
          addMessage(`SEND: ${inputValue}`);
          if (wsRef.current) {
            wsRef.current.send(inputValue);
          }
        };

        wsRef.current.onclose = (event) => {
          console.log("WebSocket connection closed:", event.code, event.reason);
          setIsConnected(false);
          addMessage(
            `CLOSE (Code: ${event.code}${
              event.reason ? `, Reason: ${event.reason}` : ""
            })`
          );
          wsRef.current = null;
        };

        wsRef.current.onmessage = (evt) => {
          console.log("Received message:", evt.data);
          addMessage(`RESPONSE: ${evt.data}`);
        };

        wsRef.current.onerror = (error) => {
          console.error("WebSocket error:", error);
          addMessage(
            `ERROR: Connection failed - check browser console for details`
          );
        };
      } catch (error: any) {
        console.error("WebSocket connection error:", error);
        addMessage(
          `ERROR: Failed to connect to WebSocket server - ${error.message}`
        );
      }
      return;
    }
    addMessage(`SEND: ${inputValue}`);
    wsRef.current.send(inputValue);
  };

  useEffect(() => {
    if (scrollAreaRef.current) {
      scrollAreaRef.current.scrollTop = scrollAreaRef.current.scrollHeight;
    }
  }, [messages]);

  return (
    <div className="container mx-auto p-6 max-w-4xl">
      <Card>
        <CardHeader>
          <CardTitle>End-to-end Test Generation</CardTitle>
          <CardDescription>
            Generate end-to-end tests that are passing with just a prompt and
            the website URL.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-4">
            <div className="space-y-4">
              <div className="flex gap-2"></div>
              <div className="flex gap-2">
                <Textarea
                  value={inputValue}
                  onChange={(e) => setInputValue(e.target.value)}
                  placeholder="Write a the website's URL, yeah, that's it."
                />
                <div className="flex flex-col gap-2">
                  {" "}
                  <Button onClick={handleSend} variant="secondary">
                    Send
                  </Button>
                  <Button
                    onClick={handleClose}
                    disabled={!isConnected}
                    variant="destructive"
                  >
                    Close
                  </Button>
                </div>
              </div>
            </div>
            <div
              ref={scrollAreaRef}
              className="h-[200px] overflow-y-auto border rounded-md p-4 bg-muted/10"
            >
              {messages.map((message, index) => (
                <div key={index} className="py-1 font-mono text-sm">
                  {message}
                </div>
              ))}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
