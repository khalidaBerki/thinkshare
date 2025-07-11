import 'package:dio/dio.dart';
import '../config/api_config.dart';
import '../features/auth/presentation/providers/auth_provider.dart';

class PaymentService {
  final Dio _dio = Dio();
  final AuthProvider _authProvider = AuthProvider();

  PaymentService() {
    _dio.options.baseUrl = ApiConfig.baseUrl;
  }

  /// Créer une session de paiement Stripe pour s'abonner à un créateur
  Future<String> createSubscriptionSession({
    required int creatorId,
    required String type, // 'paid'
  }) async {
    try {
      final token = await _authProvider.getToken();
      
      final response = await _dio.post(
        '/api/subscribe/paid',
        data: {
          'creator_id': creatorId,
          'type': type,
        },
        options: Options(
          headers: {
            'Authorization': 'Bearer $token',
            'Content-Type': 'application/json',
          },
        ),
      );

      // Le backend retourne une URL de checkout Stripe
      return response.data['checkout_url'];
    } catch (e) {
      throw Exception('Erreur lors de la création de la session de paiement: $e');
    }
  }

  /// Vérifier le statut d'abonnement d'un utilisateur
  Future<Map<String, dynamic>> getSubscriptionStatus(int creatorId) async {
    try {
      final token = await _authProvider.getToken();
      
      final response = await _dio.get(
        '/api/subscription/status/$creatorId',
        options: Options(
          headers: {
            'Authorization': 'Bearer $token',
          },
        ),
      );

      return response.data;
    } catch (e) {
      throw Exception('Erreur lors de la vérification de l\'abonnement: $e');
    }
  }
}
