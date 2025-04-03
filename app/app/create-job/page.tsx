"use client";

import { api } from "@/lib/api";
import { useMutation } from "@tanstack/react-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Textarea } from "@/components/ui/textarea";
import { toast } from "sonner";

const createJobSchema = z.object({
  prompt: z.string().min(10).max(1000),
});

const CreateJobDemo: React.FC = () => {
  const form = useForm<z.infer<typeof createJobSchema>>({
    resolver: zodResolver(createJobSchema),
    defaultValues: {
      prompt: "",
    },
  });

  const mutation = useMutation({
    mutationFn: api.createJob,
    onSettled: () => form.reset(),
    onError: (error) => toast.error(error.message),
  });

  function onSubmit(values: z.infer<typeof createJobSchema>) {
    mutation.mutate({
      prompt: values.prompt,
    });
  }

  return (
    <div className="space-y-12">
      <h2 className="mb-6 scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight transition-colors">
        Create Job
      </h2>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <FormField
            control={form.control}
            name="prompt"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Prompt</FormLabel>
                <FormControl>
                  <Textarea
                    placeholder="Tell us about the person you need to hire..."
                    className="resize-y"
                    {...field}
                  />
                </FormControl>
                <FormDescription>
                  Describe the person you need to hire
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button type="submit">
            {mutation.isPending && (
              <svg
                className="-ml-1 size-5 animate-spin text-white"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                ></circle>
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
            )}
            Submit
          </Button>
        </form>
      </Form>
      {mutation.isSuccess && mutation.data && (
        <div className="space-y-12">
          <h2 className="mb-6 scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight transition-colors">
            {mutation.data.title}
          </h2>
          <div className="space-y-4">
            <div>
              <h3 className="text-xl font-medium">Description</h3>
              <p className="mt-2 text-gray-700">{mutation.data.description}</p>
            </div>
            <div>
              <h3 className="text-xl font-medium">Requirements</h3>
              <ul className="mt-2 list-disc pl-5 space-y-1">
                {mutation.data.requirements.map((requirement, index) => (
                  <li key={index} className="text-gray-700">
                    {requirement}
                  </li>
                ))}
              </ul>
            </div>
            <div>
              <h3 className="text-xl font-medium">Responsibilities</h3>
              <ul className="mt-2 list-disc pl-5 space-y-1">
                {mutation.data.responsibilities.map((responsibility, index) => (
                  <li key={index} className="text-gray-700">
                    {responsibility}
                  </li>
                ))}
              </ul>
            </div>
            <div className="flex gap-4">
              <div>
                <h3 className="text-xl font-medium">Experience Level</h3>
                <div className="mt-2 flex items-center">
                  <div className="h-2.5 w-full rounded-full bg-gray-200">
                    <div
                      className="h-2.5 rounded-full bg-blue-600"
                      style={{
                        width: `${Math.min(100, mutation.data.experienceLevel * 20)}%`,
                      }}
                    ></div>
                  </div>
                  <span className="ml-2 text-sm font-medium text-gray-700">
                    {mutation.data.experienceLevel}/5
                  </span>
                </div>
              </div>
            </div>
            <div>
              <h3 className="text-xl font-medium">Skills</h3>
              <div className="mt-2 flex flex-wrap gap-2">
                {mutation.data.skills.map((skill, index) => (
                  <span
                    key={index}
                    className="rounded-full bg-blue-100 px-3 py-1 text-sm font-medium text-blue-800"
                  >
                    {skill}
                  </span>
                ))}
              </div>
            </div>
            <div>
              <h3 className="text-xl font-medium">Keywords</h3>
              <div className="mt-2 flex flex-wrap gap-2">
                {mutation.data.keywords.map((keyword, index) => (
                  <span
                    key={index}
                    className="rounded-full bg-gray-100 px-3 py-1 text-sm font-medium text-gray-800"
                  >
                    {keyword}
                  </span>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default function Home() {
  return <CreateJobDemo />;
}
