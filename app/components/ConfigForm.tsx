"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Eye, EyeOff } from "lucide-react";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { useUpdateConfig } from "@/lib/api";
import { fetchGet } from "@/lib/api.helpers";

// Define the config state type directly in this file
interface ConfigState {
  anthropicApiKey: string;
  sentryApiKey: string;
  techSpecification: string;
  productSpecification: string;
}

export function ConfigForm() {
  const [showAnthropicKey, setShowAnthropicKey] = useState(false);
  const [showSentryKey, setShowSentryKey] = useState(false);
  const updateConfigMutation = useUpdateConfig();
  
  // Initialize form with async defaultValues
  const form = useForm<ConfigState>({
    defaultValues: async () => {
      try {
        // Use the fetchGet helper directly
        const apiConfig = await fetchGet<ConfigState>("config");
        
        // Return the fetched config with empty string fallbacks
        return {
          anthropicApiKey: apiConfig?.anthropicApiKey || "",
          sentryApiKey: apiConfig?.sentryApiKey || "",
          techSpecification: apiConfig?.techSpecification || "",
          productSpecification: apiConfig?.productSpecification || ""
        };
      } catch (error) {
        console.error('Error fetching config:', error);
        // Return empty defaults if fetch fails
        return {
          anthropicApiKey: "",
          sentryApiKey: "",
          techSpecification: "",
          productSpecification: ""
        };
      }
    }
  });
  
  // Get form state for loading and submitting status
  const { isLoading, isSubmitting } = form.formState;

  const onSubmit = async (data: ConfigState) => {
    try {
      // Update the config in the API
      await updateConfigMutation.mutateAsync(data);
      toast.success("Configuration saved successfully");
    } catch (error) {
      toast.error("Failed to save configuration");
      console.error(error);
    }
  };

  return (
    <div className="w-full max-w-5xl mx-auto">
      {isLoading ? (
        <div className="space-y-4">
          <div className="h-4 bg-gray-200 rounded w-1/4 animate-pulse"></div>
          <div className="h-10 bg-gray-200 rounded w-full animate-pulse"></div>
          <div className="h-4 bg-gray-200 rounded w-3/4 animate-pulse"></div>
          <div className="h-4 bg-gray-200 rounded w-1/4 animate-pulse mt-6"></div>
          <div className="h-10 bg-gray-200 rounded w-full animate-pulse"></div>
          <div className="h-4 bg-gray-200 rounded w-3/4 animate-pulse"></div>
        </div>
      ) : (
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            <div className="space-y-4">
              <FormField
                control={form.control}
                name="anthropicApiKey"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Anthropic API Key</FormLabel>
                    <FormControl>
                      <div className="relative">
                        <Input
                          placeholder="sk-ant-..."
                          type={showAnthropicKey ? "text" : "password"}
                          autoComplete="off"
                          className="pr-10"
                          {...field}
                        />
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          className="absolute right-0 top-0 h-full px-3"
                          onClick={() => setShowAnthropicKey(!showAnthropicKey)}
                        >
                          {showAnthropicKey ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                          <span className="sr-only">{showAnthropicKey ? "Hide" : "Show"} API key</span>
                        </Button>
                      </div>
                    </FormControl>
                    <FormDescription>
                      Your Anthropic API key for Claude model access
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="sentryApiKey"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Sentry API Key</FormLabel>
                    <FormControl>
                      <div className="relative">
                        <Input
                          placeholder="Enter your Sentry API key"
                          type={showSentryKey ? "text" : "password"}
                          autoComplete="off"
                          className="pr-10"
                          {...field}
                        />
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          className="absolute right-0 top-0 h-full px-3"
                          onClick={() => setShowSentryKey(!showSentryKey)}
                        >
                          {showSentryKey ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                          <span className="sr-only">{showSentryKey ? "Hide" : "Show"} API key</span>
                        </Button>
                      </div>
                    </FormControl>
                    <FormDescription>
                      Your Sentry API key for error tracking integration
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

<FormField
                control={form.control}
                name="productSpecification"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Product Specification</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder="Describe your product, its features, and user requirements..."
                        className="min-h-32"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      Provide details about your product's purpose, features, and user requirements
                      that should be considered when generating tests
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="techSpecification"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Tech Specification</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder="Describe the tech stack and specifications of your application..."
                        className="min-h-32"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      Provide details about your application's tech stack, framework,
                      and any specific testing requirements
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="flex justify-end items-center">
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? "Saving..." : "Save Configuration"}
              </Button>
            </div>
      </form>
    </Form>
    )}
    </div>
  );
}
