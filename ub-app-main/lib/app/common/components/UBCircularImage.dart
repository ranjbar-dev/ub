import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:get/utils.dart';
import 'package:transparent_image/transparent_image.dart';
import '../../../generated/colors.gen.dart';

class UBCircularImage extends StatelessWidget {
  final double size;
  final String imageAddress;
  final EdgeInsets padding;
  final BoxFit fit;
  const UBCircularImage(
      {Key key,
      this.size = 33,
      this.imageAddress,
      this.padding = const EdgeInsets.all(4),
      this.fit})
      : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Container(
      width: size,
      height: size,
      padding: padding,
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(100.0),
          border: Border.all(
            width: 1,
            color: ColorName.grey36,
          ),
        ),
        padding: EdgeInsets.all(1),
        child: !(GetPlatform.isWeb)
            ? CachedNetworkImage(
                fit: fit,
                imageUrl: imageAddress,
                //imageBuilder: (context, imageProvider) => Container(
                //  decoration: BoxDecoration(
                //    image: DecorationImage(
                //        image: imageProvider,
                //        fit: BoxFit.cover,
                //        colorFilter:
                //            ColorFilter.mode(Colors.red, BlendMode.colorBurn)),
                //  ),
                //),
                progressIndicatorBuilder: (context, url, downloadProgress) =>
                    SizedBox(
                  width: 20,
                  child: CircularProgressIndicator(
                    value: downloadProgress.progress,
                    strokeWidth: 1,
                  ),
                ),
                errorWidget: (context, url, error) => Icon(Icons.error),
              )
            : FadeInImage.memoryNetwork(
                placeholder: kTransparentImage,
                image: imageAddress,
                fadeInDuration: const Duration(
                  milliseconds: 300,
                ),
              ),
      ),
    );
  }
}
