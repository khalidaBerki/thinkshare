import 'package:flutter/material.dart';
import 'package:file_picker/file_picker.dart';
import 'package:dio/dio.dart';
import 'dart:html' as html;
//import 'package:go_router/go_router.dart';
import '../../data/home_repository.dart';

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
      setState(() => _error = "Erreur lors de la sélection d'images.");
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
      setState(() => _error = "Erreur lors de la sélection de la vidéo.");
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
      setState(() => _error = "Erreur lors de la sélection des documents.");
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
        images = [
          for (final img in _images)
            if (img.bytes != null)
              MultipartFile.fromBytes(img.bytes!, filename: img.name),
        ];
      }
      if (_video != null && _video!.bytes != null) {
        video = MultipartFile.fromBytes(_video!.bytes!, filename: _video!.name);
      }
      if (_documents.isNotEmpty) {
        documents = [
          for (final doc in _documents)
            if (doc.bytes != null)
              MultipartFile.fromBytes(doc.bytes!, filename: doc.name),
        ];
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
            content: Text("Post created successfully !"),
            backgroundColor: Colors.green,
            duration: Duration(seconds: 3),
          ),
        );
        await Future.delayed(const Duration(seconds: 3));
        // Force le rechargement de la page (Flutter Web uniquement)
        // ignore: undefined_prefixed_name
        html.window.location.reload();
      }
    } on DioException catch (e) {
      setState(() {
        _error = e.response?.data['error']?.toString() ?? e.message;
      });
    } catch (e) {
      setState(() {
        _error = "Erreur inattendue : $e";
      });
    } finally {
      setState(() => _isLoading = false);
    }
  }

  Widget _mediaPreview() {
    if (_images.isNotEmpty) {
      return Wrap(
        spacing: 8,
        children: _images
            .map(
              (img) => Stack(
                alignment: Alignment.topRight,
                children: [
                  ClipRRect(
                    borderRadius: BorderRadius.circular(12),
                    child: Image.memory(
                      img.bytes!,
                      width: 80,
                      height: 80,
                      fit: BoxFit.cover,
                    ),
                  ),
                  IconButton(
                    icon: const Icon(Icons.close, size: 18, color: Colors.red),
                    onPressed: () {
                      setState(() => _images.remove(img));
                    },
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
      appBar: AppBar(title: const Text("Add Post"), centerTitle: true),
      backgroundColor: colorScheme.surface,
      body: Center(
        child: ConstrainedBox(
          constraints: BoxConstraints(maxWidth: isWide ? 600 : double.infinity),
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(18),
            child: Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Barre d'édition simple
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 6,
                    ),
                    decoration: BoxDecoration(
                      color: colorScheme.surface,
                      borderRadius: BorderRadius.circular(10),
                      border: Border.all(
                        color: colorScheme.primary.withOpacity(0.15),
                      ),
                    ),
                    child: Row(
                      children: [
                        Icon(Icons.format_bold, color: colorScheme.primary),
                        const SizedBox(width: 8),
                        Icon(Icons.format_italic, color: colorScheme.primary),
                        const SizedBox(width: 8),
                        Icon(
                          Icons.format_underline,
                          color: colorScheme.primary,
                        ),
                        const SizedBox(width: 8),
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
                  const SizedBox(height: 12),
                  TextFormField(
                    controller: _contentController,
                    minLines: 3,
                    maxLines: 8,
                    style: const TextStyle(fontFamily: 'Montserrat'),
                    decoration: InputDecoration(
                      labelText: "Content *",
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                      fillColor: colorScheme.surface,
                      filled: true,
                    ),
                    validator: (v) => v == null || v.trim().isEmpty
                        ? "Le contenu est obligatoire"
                        : null,
                  ),
                  const SizedBox(height: 16),
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
                  const SizedBox(height: 10),
                  DropdownButtonFormField<String>(
                    value: _documentType,
                    decoration: InputDecoration(
                      labelText: "Document type (optionnel)",
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    items: [
                      const DropdownMenuItem(
                        value: "report",
                        child: Text("Report"),
                      ),
                      const DropdownMenuItem(
                        value: "note",
                        child: Text("Note"),
                      ),
                      const DropdownMenuItem(
                        value: "memo",
                        child: Text("Memo"),
                      ),
                      // Ajoute d'autres types si besoin
                    ],
                    onChanged: (v) => setState(() => _documentType = v),
                  ),
                  const SizedBox(height: 16),
                  Wrap(
                    spacing: 8,
                    runSpacing: 8,
                    children: [
                      ElevatedButton.icon(
                        onPressed: _pickImages,
                        icon: const Icon(Icons.image),
                        label: const Text("Images"),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: colorScheme.primary,
                          foregroundColor: colorScheme.onPrimary,
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(8),
                          ),
                        ),
                      ),
                      ElevatedButton.icon(
                        onPressed: _pickVideo,
                        icon: const Icon(Icons.videocam),
                        label: const Text("Vidéo"),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.deepPurple,
                          foregroundColor: Colors.white,
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(8),
                          ),
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
                            borderRadius: BorderRadius.circular(8),
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  _mediaPreview(),
                  if (_error != null)
                    Padding(
                      padding: const EdgeInsets.only(top: 12),
                      child: Text(
                        _error!,
                        style: const TextStyle(
                          color: Colors.red,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                  const SizedBox(height: 24),
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      onPressed: _isLoading ? null : _submit,
                      style: ElevatedButton.styleFrom(
                        backgroundColor: colorScheme.primary,
                        foregroundColor: colorScheme.onPrimary,
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                        padding: const EdgeInsets.symmetric(vertical: 16),
                      ),
                      child: _isLoading
                          ? const CircularProgressIndicator(color: Colors.white)
                          : const Text(
                              "Submit Post",
                              style: TextStyle(
                                fontFamily: 'Montserrat',
                                fontWeight: FontWeight.bold,
                                fontSize: 17,
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
    );
  }
}
