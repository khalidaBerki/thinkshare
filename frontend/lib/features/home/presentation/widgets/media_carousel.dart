import 'package:flutter/material.dart';
import '../../../../config/api_config.dart';

class MediaCarousel extends StatelessWidget {
  final List<String> mediaUrls;

  const MediaCarousel({super.key, required this.mediaUrls});

  @override
  Widget build(BuildContext context) {
    if (mediaUrls.isEmpty) return const SizedBox.shrink();

    return SizedBox(
      height: 220,
      child: PageView.builder(
        itemCount: mediaUrls.length,
        controller: PageController(viewportFraction: 0.92),
        itemBuilder: (context, index) {
          final url = mediaUrls[index].replaceAll('\\', '/');
          final ext = url.split('.').last.toLowerCase();
          final fullUrl = '${ApiConfig.baseUrl}$url';

          if (['png', 'jpg', 'jpeg', 'gif', 'webp'].contains(ext)) {
            // Image
            return Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8.0),
              child: ClipRRect(
                borderRadius: BorderRadius.circular(18),
                child: Image.network(
                  fullUrl,
                  fit: BoxFit.cover,
                  width: double.infinity,
                  errorBuilder: (context, error, stackTrace) {
                    return Container(
                      color: Colors.grey[300],
                      child: const Center(child: Icon(Icons.broken_image)),
                    );
                  },
                ),
              ),
            );
          } else if (['mp4', 'mov', 'avi', 'webm'].contains(ext)) {
            // Vidéo (affiche une icône ou une preview si tu veux)
            return Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8.0),
              child: Container(
                decoration: BoxDecoration(
                  color: Colors.black12,
                  borderRadius: BorderRadius.circular(18),
                ),
                child: Center(
                  child: Icon(Icons.play_circle_fill, size: 64, color: Colors.deepPurple),
                ),
              ),
            );
          } else {
            // Document
            return Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8.0),
              child: Container(
                decoration: BoxDecoration(
                  color: Colors.blueGrey[50],
                  borderRadius: BorderRadius.circular(18),
                ),
                child: Center(
                  child: Icon(Icons.insert_drive_file, size: 48, color: Colors.blueGrey),
                ),
              ),
            );
          }
        },
      ),
    );
  }
}
