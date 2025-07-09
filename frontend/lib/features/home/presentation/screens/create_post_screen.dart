import 'package:flutter/material.dart';
import 'package:file_picker/file_picker.dart';
import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import '../../data/home_repository.dart';
import 'web_reload_stub.dart' if (dart.library.html) 'web_reload.dart';

class CreatePostScreen extends StatefulWidget {
  const CreatePostScreen({super.key});

  @override
  State<CreatePostScreen> createState() => _CreatePostScreenState();
}

class _CreatePostScreenState extends State<CreatePostScreen> {
  final _formKey = GlobalKey<FormState>();
  final _contentController = TextEditingController();
  String _visibility = 'public';
  String? _documentType;
  List<PlatformFile> _images = [];
  PlatformFile? _video;
  List<PlatformFile> _documents = [];
  bool _isLoading = false;
  String? _error;

  final HomeRepository _repository = HomeRepository();

  Future<void> _pickImages() async {
    try {
      final result = await FilePicker.platform.pickFiles(
        type: FileType.image,
        allowMultiple: true,
        withData: true,
      );
      if (result != null) {
        setState(() {
          _images = result.files.take(10).toList();
          _video = null;
          _documents = [];
        });
      }
    } catch (e) {
      setState(() => _error = "Error while picking images.");
    }
  }

  Future<void> _pickVideo() async {
    try {
      final result = await FilePicker.platform.pickFiles(
        type: FileType.video,
        allowMultiple: false,
        withData: true,
      );
      if (result != null && result.files.isNotEmpty) {
        setState(() {
          _video = result.files.first;
          _images = [];
          _documents = [];
        });
      }
    } catch (e) {
      setState(() => _error = "Error while picking video.");
    }
  }

  Future<void> _pickDocuments() async {
    try {
      final result = await FilePicker.platform.pickFiles(
        type: FileType.custom,
        allowedExtensions: [
          'pdf',
          'doc',
          'docx',
          'ppt',
          'pptx',
          'xls',
          'xlsx',
          'txt',
        ],
        allowMultiple: true,
        withData: true,
      );
      if (result != null) {
        setState(() {
          _documents = result.files.take(5).toList();
          _images = [];
          _video = null;
        });
      }
    } catch (e) {
      setState(() => _error = "Error while picking documents.");
    }
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      List<MultipartFile>? images;
      MultipartFile? video;
      List<MultipartFile>? documents;

      if (_images.isNotEmpty) {
        images = _images
            .where((img) => img.bytes != null)
            .map((img) => MultipartFile.fromBytes(img.bytes!, filename: img.name))
            .toList();
      }
      if (_video != null && _video!.bytes != null) {
        video = MultipartFile.fromBytes(_video!.bytes!, filename: _video!.name);
      }
      if (_documents.isNotEmpty) {
        documents = _documents
            .where((doc) => doc.bytes != null)
            .map((doc) => MultipartFile.fromBytes(doc.bytes!, filename: doc.name))
            .toList();
      }

      await _repository.createPost(
        content: _contentController.text.trim(),
        visibility: _visibility,
        documentType: _documentType,
        images: images,
        video: video,
        documents: documents,
      );

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text("Post created successfully!"),
            backgroundColor: Colors.green,
            duration: Duration(seconds: 3),
          ),
        );
        await Future.delayed(const Duration(seconds: 3));
        if (kIsWeb) {
          reloadPage();
        } else {
          if (Navigator.canPop(context)) {
            Navigator.pop(context);
          }
        }
      }
    } on DioException catch (e) {
      setState(() {
        _error = e.response?.data['error']?.toString() ?? e.message;
      });
    } catch (e) {
      setState(() {
        _error = "Unexpected error: $e";
      });
    } finally {
      setState(() => _isLoading = false);
    }
  }

  Widget _mediaPreview() {
    if (_images.isNotEmpty) {
      return Wrap(
        spacing: 10,
        runSpacing: 10,
        children: _images
            .map(
              (img) => Stack(
                alignment: Alignment.topRight,
                children: [
                  ClipRRect(
                    borderRadius: BorderRadius.circular(16),
                    child: Image.memory(
                      img.bytes!,
                      width: 90,
                      height: 90,
                      fit: BoxFit.cover,
                    ),
                  ),
                  Material(
                    color: Colors.transparent,
                    child: InkWell(
                      borderRadius: BorderRadius.circular(16),
                      onTap: () => setState(() => _images.remove(img)),
                      child: const Padding(
                        padding: EdgeInsets.all(4),
                        child: Icon(Icons.close, size: 18, color: Colors.red),
                      ),
                    ),
                  ),
                ],
              ),
            )
            .toList(),
      );
    }
    if (_video != null) {
      return ListTile(
        leading: const Icon(Icons.videocam, color: Colors.deepPurple),
        title: Text(_video!.name),
        trailing: IconButton(
          icon: const Icon(Icons.close, color: Colors.red),
          onPressed: () => setState(() => _video = null),
        ),
      );
    }
    if (_documents.isNotEmpty) {
      return Wrap(
        spacing: 8,
        children: _documents
            .map(
              (doc) => Chip(
                label: Text(doc.name),
                onDeleted: () => setState(() => _documents.remove(doc)),
                backgroundColor: Colors.blueGrey.shade50,
                labelStyle: const TextStyle(fontWeight: FontWeight.w500),
              ),
            )
            .toList(),
      );
    }
    return const SizedBox.shrink();
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final isWide = MediaQuery.of(context).size.width > 600;

    return Scaffold(
      appBar: AppBar(
        title: const Text(
          "Create Post",
          style: TextStyle(
            fontFamily: 'Montserrat',
            fontWeight: FontWeight.bold,
            fontSize: 20,
            letterSpacing: 0.2,
          ),
        ),
        centerTitle: true,
        elevation: 2,
        backgroundColor: colorScheme.surface,
        shadowColor: colorScheme.primary.withOpacity(0.08),
        surfaceTintColor: colorScheme.primary,
        actions: const [],
      ),
      backgroundColor: colorScheme.surface,
      body: Center(
        child: ConstrainedBox(
          constraints: BoxConstraints(maxWidth: isWide ? 600 : double.infinity),
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(22),
            child: Material(
              elevation: 3,
              borderRadius: BorderRadius.circular(24),
              color: colorScheme.surface,
              child: Padding(
                padding: const EdgeInsets.all(22),
                child: Form(
                  key: _formKey,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // Simple edit bar (icons only for style)
                      Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 8,
                        ),
                        decoration: BoxDecoration(
                          color: colorScheme.surfaceContainerHighest.withOpacity(0.7),
                          borderRadius: BorderRadius.circular(14),
                        ),
                        child: Row(
                          children: [
                            Icon(Icons.format_bold, color: colorScheme.primary),
                            const SizedBox(width: 10),
                            Icon(
                              Icons.format_italic,
                              color: colorScheme.primary,
                            ),
                            const SizedBox(width: 10),
                            Icon(
                              Icons.format_underline,
                              color: colorScheme.primary,
                            ),
                            const SizedBox(width: 10),
                            Icon(
                              Icons.format_align_left,
                              color: colorScheme.primary,
                            ),
                            Icon(
                              Icons.format_align_center,
                              color: colorScheme.primary,
                            ),
                            Icon(
                              Icons.format_align_right,
                              color: colorScheme.primary,
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 16),
                      TextFormField(
                        controller: _contentController,
                        minLines: 3,
                        maxLines: 8,
                        style: const TextStyle(fontFamily: 'Montserrat'),
                        decoration: InputDecoration(
                          labelText: "Content *",
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(16),
                          ),
                          fillColor: colorScheme.surfaceContainerHighest.withOpacity(
                            0.5,
                          ),
                          filled: true,
                        ),
                        validator: (v) => v == null || v.trim().isEmpty
                            ? "Content is required"
                            : null,
                      ),
                      const SizedBox(height: 18),
                      Row(
                        children: [
                          Switch(
                            value: _visibility == 'public',
                            onChanged: (v) => setState(
                              () => _visibility = v ? 'public' : 'private',
                            ),
                            activeColor: colorScheme.primary,
                          ),
                          Text(
                            _visibility == 'public'
                                ? "Public Post"
                                : "Private Post",
                            style: TextStyle(
                              color: _visibility == 'public'
                                  ? colorScheme.primary
                                  : colorScheme.error,
                              fontWeight: FontWeight.bold,
                              fontFamily: 'Montserrat',
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 12),
                      DropdownButtonFormField<String>(
                        value: _documentType,
                        decoration: InputDecoration(
                          labelText: "Document type (optional)",
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(14),
                          ),
                        ),
                        items: const [
                          DropdownMenuItem(
                            value: "report",
                            child: Text("Report"),
                          ),
                          DropdownMenuItem(value: "note", child: Text("Note")),
                          DropdownMenuItem(value: "memo", child: Text("Memo")),
                        ],
                        onChanged: (v) => setState(() => _documentType = v),
                      ),
                      const SizedBox(height: 18),
                      Wrap(
                        spacing: 10,
                        runSpacing: 10,
                        children: [
                          ElevatedButton.icon(
                            onPressed: _pickImages,
                            icon: const Icon(Icons.image),
                            label: const Text("Images"),
                            style: ElevatedButton.styleFrom(
                              backgroundColor: colorScheme.primary,
                              foregroundColor: colorScheme.onPrimary,
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(10),
                              ),
                              elevation: 2,
                            ),
                          ),
                          ElevatedButton.icon(
                            onPressed: _pickVideo,
                            icon: const Icon(Icons.videocam),
                            label: const Text("Video"),
                            style: ElevatedButton.styleFrom(
                              backgroundColor: Colors.deepPurple,
                              foregroundColor: Colors.white,
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(10),
                              ),
                              elevation: 2,
                            ),
                          ),
                          ElevatedButton.icon(
                            onPressed: _pickDocuments,
                            icon: const Icon(Icons.insert_drive_file),
                            label: const Text("Documents"),
                            style: ElevatedButton.styleFrom(
                              backgroundColor: Colors.blueGrey,
                              foregroundColor: Colors.white,
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(10),
                              ),
                              elevation: 2,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 14),
                      _mediaPreview(),
                      if (_error != null)
                        Padding(
                          padding: const EdgeInsets.only(top: 14),
                          child: Text(
                            _error!,
                            style: const TextStyle(
                              color: Colors.red,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                      const SizedBox(height: 28),
                      SizedBox(
                        width: double.infinity,
                        child: ElevatedButton(
                          onPressed: _isLoading ? null : _submit,
                          style: ElevatedButton.styleFrom(
                            backgroundColor: colorScheme.primary,
                            foregroundColor: colorScheme.onPrimary,
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(16),
                            ),
                            padding: const EdgeInsets.symmetric(vertical: 18),
                            elevation: 3,
                          ),
                          child: _isLoading
                              ? const CircularProgressIndicator(
                                  color: Colors.white,
                                )
                              : const Text(
                                  "Submit Post",
                                  style: TextStyle(
                                    fontFamily: 'Montserrat',
                                    fontWeight: FontWeight.bold,
                                    fontSize: 18,
                                  ),
                                ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
