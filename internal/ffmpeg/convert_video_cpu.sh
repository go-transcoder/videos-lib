#!/bin/sh
file=$1     #full video path on storage: tenant_folder/uploads/file_name_of_video.mp4
path=$2     #full path where to convert different quality of video: tenant_folder/uploads/file_name_of_video/
ffmpeg_cv="libx264"
#TODO refactor if block

mkdir -p $path/240p/
ffmpeg -y -loglevel error -hide_banner -i $file -filter_complex "[0]split=1[s0]; \
   [s0]scale='min(trunc(oh*a/2)*2, iw)':'min(240, ih)'[v0]" \
  -map "[v0]" -map 0:a? -c:v $ffmpeg_cv -pix_fmt yuv420p -c:a aac -preset:v fast -profile:v high -b:v 500k -bf:v 3 -r 24 -g 24  $path/240p/video.mp4

