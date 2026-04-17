import Image from "next/image";
import React, { useEffect, useState } from "react";
import { RxCross2 } from "react-icons/rx";
import { BiEdit } from "react-icons/bi";
import SelectList from "@/components/ui/SelectList";
import { Controller, useFormContext } from "react-hook-form";
import { TCourseSchema } from "@/schema/course.schema";

const videoSource = [
  {
    id: 1,
    name: "Youtube",
    value: "youtube",
  },
  {
    id: 2,
    name: "Vimeo",
    value: "vimeo",
  },
];

export default function IntroVideo() {
  const [isAddVideo, setIsAddVideo] = useState(false);
  const [videoUrl, setVideoUrl] = useState("");
  const [videoId, setVideoId] = useState("");

  const {
    register,
    control,
    formState: { errors },
    resetField,
    trigger,
    watch,
  } = useFormContext<TCourseSchema>();

  const extractYouTubeVideoId = (url: string) => {
    const match = url.match(
      /(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/watch\?v=|youtu\.be\/)([^\s&]+)/
    );
    return match ? match[1] : null;
  };

  const handleOk = () => {
    trigger("intro_video.source");
    const id = extractYouTubeVideoId(videoUrl);
    if (id) {
      setVideoId(id);
      setIsAddVideo(false);
    }
  };

  const handleCancel = () => {
    setIsAddVideo(false);
    setVideoUrl("");
    setVideoId("");
    resetField("intro_video.source");
  };

  const handleEdit = () => {
    setIsAddVideo(true);
  };

  const thumbnailUrl = videoId
    ? `https://img.youtube.com/vi/${videoId}/hqdefault.jpg`
    : "";

  const handleRemoveVideo = () => {
    setVideoId("");
    setVideoUrl("");
  };

  useEffect(() => {
    if (watch("intro_video.source")) {
      setIsAddVideo(true);
    }
  }, [watch("intro_video.source")]);

  return (
    <>
      {videoId && !isAddVideo ? (
        <div className="relative group mb-4 cursor-pointer bg-white p-3 border rounded-xl">
          <div className="flex items-center justify-between mb-2">
            <p className="text-sm truncate">{videoUrl}</p>
            <div className="flex items-center gap-2 text-gray-500">
              <button onClick={handleEdit} className="p-1">
                <BiEdit />
              </button>
              <button onClick={handleRemoveVideo} className="p-1">
                <RxCross2 />
              </button>
            </div>
          </div>
          <div className="h-56 w-full">
            <Image
              src={thumbnailUrl}
              alt="YouTube thumbnail"
              width={400}
              height={400}
              className="rounded-md w-full h-full object-cover"
            />
          </div>
        </div>
      ) : null}

      {isAddVideo ? (
        <div className="bg-white p-4 border rounded-md">
          <Controller
            control={control}
            name="intro_video.type"
            defaultValue="youtube"
            render={({ field: { onChange, value } }) => (
              <SelectList
                options={videoSource}
                value={videoSource.find((s) => s.value === value)}
                onChange={(s) => onChange(s.value)}
                className="w-full"
              />
            )}
          />
          <textarea
            placeholder="Paste Video URL"
            className="outline-none w-full border border-dashed mt-5 rounded-md min-h-20 text-sm px-3 py-1.5"
            {...register("intro_video.source")}
          />
          {errors.intro_video?.source && (
            <p className="text-red-500 text-sm mt-1">
              {errors.intro_video.source.message}
            </p>
          )}
          <div className="flex items-center justify-end text-sm mt-3">
            <button type="button" onClick={handleCancel} className="p-3">
              Cancel
            </button>
            <button
              type="button"
              onClick={handleOk}
              className="bg-primary/10 text-primary px-3 py-1 rounded-md"
            >
              Ok
            </button>
          </div>
        </div>
      ) : !videoId ? (
        <div className="border border-dashed rounded-lg h-44 bg-white flex flex-col items-center justify-center cursor-pointer">
          <button
            onClick={() => setIsAddVideo(true)}
            className="text-sm my-2 bg-blue-200 text-blue-700 font-medium px-2.5 py-1 rounded-md"
          >
            + Add from URL
          </button>
        </div>
      ) : null}
    </>
  );
}
