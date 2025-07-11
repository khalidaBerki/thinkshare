import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart';
import 'package:dio/dio.dart';
import 'package:url_launcher/url_launcher.dart';
import '../../../../services/payment_service.dart';

class UpgradeBanner extends StatelessWidget {
  final int? creatorId;
  final double? monthlyPrice;
  final String? username;
  const UpgradeBanner({this.creatorId, this.monthlyPrice, this.username, super.key});

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final bool canPay = creatorId != null;
    return Container(
      margin: const EdgeInsets.symmetric(vertical: 16),
      padding: const EdgeInsets.all(18),
      decoration: BoxDecoration(
        color: colorScheme.primary.withOpacity(0.08),
        border: Border.all(color: colorScheme.primary),
        borderRadius: BorderRadius.circular(14),
      ),
      child: Column(
        children: [
          Text(
            "🔒 Ce contenu est réservé. Abonnez-vous ou payez pour y accéder !",
            style: TextStyle(
              color: colorScheme.primary,
              fontFamily: 'Montserrat',
              fontWeight: FontWeight.bold,
              fontSize: 15,
            ),
            textAlign: TextAlign.center,
          ),
          if (monthlyPrice != null && monthlyPrice! > 0)
            Padding(
              padding: const EdgeInsets.only(top: 6.0, bottom: 2.0),
              child: Text(
                'Abonnement : ${monthlyPrice!.toStringAsFixed(2)} €/mois',
                style: TextStyle(
                  color: colorScheme.primary,
                  fontWeight: FontWeight.w600,
                  fontSize: 14,
                ),
                textAlign: TextAlign.center,
              ),
            ),
          const SizedBox(height: 10),
          if (canPay)
            Column(
              children: [
                ElevatedButton.icon(
                  icon: const Icon(Icons.payment),
                  label: Text(
                    monthlyPrice != null && monthlyPrice! > 0
                      ? 'Payer ${monthlyPrice!.toStringAsFixed(2)} €/mois pour accéder'
                      : 'Accéder et payer',
                  ),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: colorScheme.primary,
                    foregroundColor: colorScheme.onPrimary,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                    padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 10),
                  ),
                  onPressed: () => _handlePayAction(context),
                ),
                if (username != null && username!.isNotEmpty)
                  Padding(
                    padding: const EdgeInsets.only(top: 6.0),
                    child: Text(
                      '@$username',
                      style: TextStyle(
                        color: colorScheme.secondary,
                        fontSize: 13,
                        fontStyle: FontStyle.italic,
                      ),
                    ),
                  ),
              ],
            )
          else
            Text(
              "Créateur inconnu ou non disponible pour le paiement.",
              style: TextStyle(color: colorScheme.error, fontSize: 14),
              textAlign: TextAlign.center,
            ),
        ],
      ),
    );
  }

  void _handlePayAction(BuildContext context) async {
    if (creatorId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Erreur: ID du créateur manquant')),
      );
      return;
    }
    try {
      final paymentService = PaymentService();
      final checkoutUrl = await paymentService.createSubscriptionSession(
        creatorId: creatorId!,
        type: 'paid',
      );
      final uri = Uri.parse(checkoutUrl);
      if (await canLaunchUrl(uri)) {
        // Use different launch modes based on platform
        if (kIsWeb) {
          // For web, use platformDefault which opens in the same tab
          await launchUrl(uri, mode: LaunchMode.platformDefault);
        } else {
          // For mobile, use externalApplication
          await launchUrl(uri, mode: LaunchMode.externalApplication);
        }
      } else {
        throw Exception('Impossible d\'ouvrir le lien Stripe');
      }
    } catch (e) {
      if (e is DioException && e.response?.data != null) {
        String backendMsg = '';
        final data = e.response?.data;
        if (data is Map && data['error'] != null) {
          backendMsg = data['error'].toString();
        } else if (data is String) {
          backendMsg = data;
        }
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Erreur Stripe: ${backendMsg.isNotEmpty ? backendMsg : e.toString()}'),
            backgroundColor: Colors.red,
          ),
        );
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erreur: $e'), backgroundColor: Colors.red),
        );
      }
    }
  }
}
